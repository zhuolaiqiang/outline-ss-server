// Copyright 2018 Jigsaw Operations LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"container/list"
	"net/netip"
	"sync"

	"github.com/Jigsaw-Code/outline-sdk/transport/shadowsocks"
)

// Don't add a tag if it would reduce the salt entropy below this amount.
const minSaltEntropy = 16

// 1. ID (string)
// 含义: 唯一标识符，通常对应于配置文件中的 Access Key ID（如 "user-1", "secret-key-0"）。
// 用途:
// 监控: 在 Metrics 中，我们用 ID 来区分流量属于哪个用户。例如 Prometheus 里的标签 access_key="user-1"。
// 日志: 出错时打印 ID，方便运维定位是谁的连接出了问题。
// 2. CryptoKey (*shadowsocks.EncryptionKey)
// 含义: 预先计算好的加密密钥对象。
// 用途:
// 包含核心的 AES/Chacha20 密钥。
// 性能优化: Shadowsocks 的密钥通常是从用户密码（字符串）通过 HKDF 或 PBKDF2 算法派生出来的。这个过程很耗 CPU。我们只在启动时计算一次，然后存其结果（CryptoKey），之后每次连接直接用，避免重复计算。
// 3. SaltGenerator (ServerSaltGenerator)
// 含义: 服务端发送数据时，生成 Salt（盐）的策略。
// 用途:
// 标准模式 (Random): 纯随机生成盐。
// 防重放模式 (Marked): 如果配置允许，Outline 会使用一种特殊的“标记盐”生成器。它生成的盐不仅是随机的，还包含一个由 Secret 派生的隐蔽 Tag。
// 反向重放保护: 这样客户端可以用同样的 Secret 验证这个盐是否真的是由合法的服务器生成的，从而防止中间人重放旧的服务器数据包。
// 4. lastClientIP (netip.Addr)
// 含义: 最近一次使用这个 Key 成功的客户端 IP。
// 用途: 性能优化 (Cache Locality)。
// 场景: 张三的手机 IP 是 1.2.3.4，他用的是 Key_A。
// 试错成本: 当服务器收到来自 1.2.3.4 的新连接时，如果不加优化，服务器得把 Key_A、Key_B、Key_C... 挨个试一遍。
// 优化逻辑: 服务器记住了 Key_A.lastClientIP = 1.2.3.4。
// 结果: 下次 1.2.3.4 再来连接，服务器会优先尝试 Key_A。如果解密成功，就省去了试其他 Key 的 CPU 开销。这就是
// SnapshotForClientIP函数的核心逻辑。
// 总结:
// CipherEntry
// 不仅仅是一个静态配置，它还是一个带有运行时状态（lastClientIP）的智能对象。它既包含了加密所需的静态材料，也包含了用于防御（SaltGenerator）和加速（lastClientIP）的动态组件。

// CipherEntry holds a Cipher with an identifier.
// The public fields are constant, but lastClientIP is mutable under cipherList.mu.
type CipherEntry struct {
	ID            string
	CryptoKey     *shadowsocks.EncryptionKey
	SaltGenerator ServerSaltGenerator
	lastClientIP  netip.Addr
}

// MakeCipherEntry constructs a CipherEntry.
func MakeCipherEntry(id string, cryptoKey *shadowsocks.EncryptionKey, secret string) CipherEntry {
	var saltGenerator ServerSaltGenerator
	if cryptoKey.SaltSize()-serverSaltMarkLen >= minSaltEntropy {
		// Mark salts with a tag for reverse replay protection.
		saltGenerator = NewServerSaltGenerator(secret)
	} else {
		// Adding a tag would leave too little randomness to protect
		// against accidental salt reuse, so don't mark the salts.
		saltGenerator = RandomServerSaltGenerator
	}
	return CipherEntry{
		ID:            id,
		CryptoKey:     cryptoKey,
		SaltGenerator: saltGenerator,
	}
}

// CipherList is a thread-safe collection of CipherEntry elements that allows for
// snapshotting and moving to front.
type CipherList interface {
	// Returns a snapshot of the cipher list optimized for this client IP
	SnapshotForClientIP(clientIP netip.Addr) []*list.Element
	MarkUsedByClientIP(e *list.Element, clientIP netip.Addr)
	// Update replaces the current contents of the CipherList with `contents`,
	// which is a List of *CipherEntry.  Update takes ownership of `contents`,
	// which must not be read or written after this call.
	Update(contents *list.List)
}

type cipherList struct {
	CipherList
	list *list.List
	mu   sync.RWMutex
}

// NewCipherList creates an empty CipherList
func NewCipherList() CipherList {
	return &cipherList{list: list.New()}
}

func matchesIP(e *list.Element, clientIP netip.Addr) bool {
	c := e.Value.(*CipherEntry)
	return clientIP != netip.Addr{} && clientIP == c.lastClientIP
}

// 调用时机:
// TCP: 每有一个新的 TCP 连接 建立时调用一次 (
// findAccessKey
// )。
// UDP: 每收到一个来自新客户端地址 (IP:Port) 的 UDP 数据包时调用一次 (
// findAccessKeyUDP
// )。
// 注意：对于 UDP，如果是同一个 ClientIP:ClientPort 发来的后续数据包，因为已经建立了 NAT 映射（
// natmap
// ），会直接查表，不会再调这个函数。只有当这是一条全新的“会话”时才会调用。
// 频率: 非常高 (Very High)。 它是握手的第一步。如果有大量的并发连接请求（或者遭到了 DoS 攻击），这个函数会被疯狂调用。

// 性能瓶颈风险:

// 它持有 cl.mu.RLock() (读锁)。虽然读锁之间不互斥，但会阻塞写锁。
// 它内部有两个 for 循环遍历整个链表。如果 Key 的数量非常多（例如几千个），且并发非常高，这里可能成为 CPU 密集型的热点。
// !!!建议：如果你真的需要支持单机几万用户，Outline 目前的这种架构可能不是最优解。
//
//	通常这种量级会采用专门的 多端口 方案（每个端口只分配部分用户），或者改用基于 Hash 的查找（但 Shadowsocks 协议本身不支持发送 UserID，所以很难 Hash）。
// 	方案 2：读写分离 (Copy-On-Write / RCU)
// 现状问题：使用了
// list
//  链表，虽然写操作极快（头部移动），但读操作（全量遍历）导致严重的 Cache Miss。

// 优化方案： 放弃链表，改为维护一个全局的、排序好的数组（Slice）。

// Default List: 维护一个默认顺序的 Slice (所有 Key)。
// IP Cache: 维护一个 Map[ClientIP] -> *CipherEntry 的缓存。
// 逻辑变更：

// Read (握手时):
// 先查 IP Cache。如果有记录，直接拿那个 Key 先试一下（O(1)）。
// 如果那个 Key 试得不对，再回落到遍历 Default List。
// 关键点：Default List 是一个连续内存的 Slice，遍历速度比链表快得多（Cache Friendly）。
// Write (成功时):
// 只更新 IP Cache。
// 不再去调整全量列表的顺序（省去了全局写锁）。
// 收益：不仅内存访问快了，而且彻底消除了
// Snapshot
//
//	的大内存分配。
func (cl *cipherList) SnapshotForClientIP(clientIP netip.Addr) []*list.Element {
	cl.mu.RLock()
	defer cl.mu.RUnlock()
	// 这里也是性能瓶颈
	cipherArray := make([]*list.Element, cl.list.Len())
	i := 0
	// First pass: put all ciphers with matching last known IP at the front.
	for e := cl.list.Front(); e != nil; e = e.Next() {
		if matchesIP(e, clientIP) {
			cipherArray[i] = e
			i++
		}
	}
	// Second pass: include all remaining ciphers in recency order.
	for e := cl.list.Front(); e != nil; e = e.Next() {
		if !matchesIP(e, clientIP) {
			cipherArray[i] = e
			i++
		}
	}
	return cipherArray
}

func (cl *cipherList) MarkUsedByClientIP(e *list.Element, clientIP netip.Addr) {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cl.list.MoveToFront(e)

	c := e.Value.(*CipherEntry)
	c.lastClientIP = clientIP
}

func (cl *cipherList) Update(src *list.List) {
	cl.mu.Lock()
	cl.list = src
	cl.mu.Unlock()
}
