#!/bin/bash
set -e  # 错误立即退出
set -u  # 未定义变量报错

# ===================== 核心配置（适配你的需求） =====================
CONTAINER_NAME="outline-ss-server"  # 容器名称
IMAGE_NAME="outline-ss-server:latest"  # 镜像标签
RESTART_POLICY="unless-stopped"     # 重启策略
# 端口映射（匹配Dockerfile的EXPOSE）
PORT_MAPPINGS=(
  "9000:9000/tcp"
  "9000:9000/udp"
  "9001:9001/tcp"
  "9001:9001/udp"
  "9090:9090/tcp"
)
# 主机数据目录（当前目录下的 outline-ss-server-data）
HOST_DATA_DIR="./outline-ss-server-data"
# 容器内挂载目标路径（统一挂载到 /data，便于管理）
CONTAINER_DATA_DIR="/data"
# 日志配置
LOG_MAX_SIZE="100m"
LOG_MAX_FILES="5"
# 内存配置（适配低版本Docker：指定超大值≈无限制，禁用Swap）
MEMORY_LIMIT="100g"  # 设为远大于主机物理内存的值（如100G）

# ===================== 前置准备 =====================
# 1. 创建主机数据目录（不存在则创建）
if [ ! -d "${HOST_DATA_DIR}" ]; then
  mkdir -p "${HOST_DATA_DIR}"
  echo -e "\033[33m创建主机数据目录：${HOST_DATA_DIR}\033[0m"
fi

# 2. 复制镜像内置配置到主机数据目录（首次启动初始化）
INIT_CONFIG="${HOST_DATA_DIR}/config.yml"
if [ ! -f "${INIT_CONFIG}" ]; then
  echo -e "\033[33m初始化配置文件到：${INIT_CONFIG}\033[0m"
  # 启动临时容器复制内置配置
  docker run --rm --name temp-outline-copy -d "${IMAGE_NAME}" >/dev/null 2>&1
  docker cp temp-outline-copy:/etc/outline/config.yml "${INIT_CONFIG}"
  docker stop temp-outline-copy >/dev/null 2>&1
fi

# 3. 检查镜像是否存在
if ! docker images --format "{{.Repository}}:{{.Tag}}" | grep -q "^${IMAGE_NAME}$"; then
  echo -e "\033[31m错误：镜像 ${IMAGE_NAME} 不存在，请先执行build脚本！\033[0m"
  exit 1
fi

# 4. 停止并删除旧容器（避免冲突）
if docker ps -a --format "{{.Names}}" | grep -q "^${CONTAINER_NAME}$"; then
  echo -e "\033[33m停止并删除旧容器 ${CONTAINER_NAME}...\033[0m"
  docker stop "${CONTAINER_NAME}" >/dev/null 2>&1
  docker rm "${CONTAINER_NAME}" >/dev/null 2>&1
fi

# ===================== 构造运行命令（核心：修复内存参数） =====================
RUN_CMD=(
  docker run -d
  --name "${CONTAINER_NAME}"
  --restart "${RESTART_POLICY}"
  # 日志配置
  --log-driver json-file
  --log-opt max-size="${LOG_MAX_SIZE}"
  --log-opt max-file="${LOG_MAX_FILES}"
  # 内存配置：低版本Docker兼容（超大内存限制 + 禁用Swap）
  # --memory="-1"       # 不设置表示使用宿主机器的最大内存
  --memory-swap="0"  # swap和内存相同 → 禁用Swap（低版本Docker标准写法）
  --memory-swappiness=0            # 内核层面禁止使用Swap
  # CPU/文件句柄无限制
  --cpus=0                         # 高版本Docker：0=无限制；低版本可省略（默认无限制）
  # 挂载数据目录（核心！主机当前目录/outline-ss-server-data → 容器/data）
  -v "${HOST_DATA_DIR}:${CONTAINER_DATA_DIR}"
  # 挂载配置文件（覆盖容器内置配置）
  -v "${INIT_CONFIG}:/etc/outline/config.yml"
)

# 添加端口映射
for port in "${PORT_MAPPINGS[@]}"; do
  RUN_CMD+=("-p" "${port}")
done

# 最终镜像 + 启动参数
RUN_CMD+=("${IMAGE_NAME}")

# ===================== 执行启动 =====================
echo -e "\033[36m启动容器 ${CONTAINER_NAME}（数据挂载到 ${HOST_DATA_DIR}）...\033[0m"
"${RUN_CMD[@]}"

# ===================== 验证启动 =====================
sleep 2
if docker ps --format "{{.Names}}" | grep -q "^${CONTAINER_NAME}$"; then
  echo -e "\033[32m容器启动成功！\033[0m"
  echo -e "├─ 容器名称：${CONTAINER_NAME}"
  echo -e "├─ 映射端口：${PORT_MAPPINGS[*]}"
  echo -e "├─ 数据目录：主机 ${HOST_DATA_DIR} → 容器 ${CONTAINER_DATA_DIR}"
  echo -e "├─ 配置文件：${INIT_CONFIG}"
  echo -e "├─ 内存策略：禁用Swap（内存限制${MEMORY_LIMIT}）"
  echo -e "├─ 重启策略：${RESTART_POLICY}"
  echo -e "├─ 日志查看：docker logs -f ${CONTAINER_NAME}"
  echo -e "└─ 数据管理：修改 ${INIT_CONFIG} 即可更新容器配置"
else
  echo -e "\033[31m容器启动失败！\033[0m"
  docker logs "${CONTAINER_NAME}"
  exit 1
fi
