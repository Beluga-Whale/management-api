# ใช้ Go เวอร์ชันล่าสุด
FROM golang:1.22-alpine

# ติดตั้ง air (hot reload)
RUN apk add --no-cache git && go install github.com/cosmtrek/air@latest

# ตั้งค่า Working Directory
WORKDIR /app

# คัดลอกโค้ดทั้งหมด
COPY . .

# เปิดพอร์ต 8080
EXPOSE 8080

# ใช้ air reload สำหรับ Dev Mode
CMD ["air"]
