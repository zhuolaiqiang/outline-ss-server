#!/bin/bash
set -e  # 遇到错误立即退出
set -x  # 显示执行的命令（调试用）

# ===================== 新增：创建images目录 =====================
# 确保当前目录下的images文件夹存在，不存在则创建
IMAGES_DIR="./images"
if [ ! -d "${IMAGES_DIR}" ]; then
  mkdir -p "${IMAGES_DIR}"
  echo -e "\033[33m创建目录：${IMAGES_DIR}\033[0m"
fi

# 动态获取Git版本号（替代固定v1.0）
VERSION=$(git rev-parse --short HEAD 2>/dev/null || echo "dev")

# ===================== 新增：生成精确到秒的时间戳 =====================
# 时间格式：YYYYMMDD_HHMMSS（年-月-日_时-分-秒），避免空格/特殊字符
BUILD_TIME=$(date +"%Y-%m-%d+%H-%M-%S")

# 构建镜像
docker build \
  -t outline-ss-server:latest \
  -t outline-ss-server:"${VERSION}" \
  --build-arg VERSION="${VERSION}" \
  --progress=plain \
  -f Dockerfile \
  .

# ===================== 新增：导出镜像到images文件夹（含时间戳） =====================
# 定义镜像导出文件名（Git版本 + 精确时间戳）
IMAGE_FILE="${IMAGES_DIR}/outline-ss-server_${VERSION}_${BUILD_TIME}.tar"

# 导出镜像（推荐用docker save，保留所有标签和层）
echo -e "\033[33m开始导出镜像到：${IMAGE_FILE}\033[0m"
docker save outline-ss-server:latest outline-ss-server:"${VERSION}" -o "${IMAGE_FILE}"

# 可选：压缩镜像文件（节省空间）
# gzip -9 "${IMAGE_FILE}"
# IMAGE_FILE="${IMAGE_FILE}.gz"

# 构建+导出成功提示
echo -e "\033[32m镜像构建并导出成功！\033[0m"
echo -e "\033[36m镜像文件路径：${IMAGE_FILE}\033[0m"
docker images outline-ss-server | head -2
ls -lh "${IMAGES_DIR}"  # 显示导出的镜像文件大小
