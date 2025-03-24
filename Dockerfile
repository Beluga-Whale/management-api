# ใช้ Go เวอร์ชันล่าสุด
FROM golang:1.22-alpine AS builder

# ตั้งค่า Working Directory
WORKDIR /app

# คัดลอก go.mod และ go.sum เพื่อติดตั้ง dependencies ก่อน
COPY go.mod go.sum ./
RUN go mod download

# คัดลอกโค้ดทั้งหมดเข้าไปใน Container
COPY . .

# สร้างไฟล์ Binary ที่ชื่อ main
RUN go build -o main .

# ใช้ Alpine เพื่อลดขนาดของ Container
FROM alpine:latest

WORKDIR /app

# คัดลอกไฟล์ Binary ที่ build เสร็จแล้ว
COPY --from=builder /app/main .

# เปิดพอร์ต 8080
EXPOSE 8080

# รันแอป Go
CMD ["./main"]
