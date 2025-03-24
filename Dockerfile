FROM golang:1.22-alpine

# ติดตั้ง air (hot reload)
RUN apk add --no-cache git && go install github.com/cosmtrek/air@latest

# ตั้งค่า Working Directory
WORKDIR /app

# คัดลอกโค้ดทั้งหมด
COPY . .

# โหลดค่า env (Docker Compose จะจัดการให้)
ENV PORT_API=:8080

# เปิดพอร์ต API
EXPOSE 8080

# ใช้ air reload สำหรับ Dev Mode
CMD ["air"]
