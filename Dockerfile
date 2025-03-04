# ใช้ Golang base image
FROM golang:1.21 as builder

# กำหนด working directory ภายใน container
WORKDIR /app

# คัดลอกไฟล์ go.mod และ go.sum ไปยัง container
COPY go.mod go.sum ./

# ดาวน์โหลด dependencies
RUN go mod download

# คัดลอกไฟล์ทั้งหมดของโปรเจกต์ไปยัง container
COPY . .

# สร้าง binary file
RUN go build -o app cmd/main.go

# ใช้ lightweight base image สำหรับ production
FROM alpine:latest  

# ตั้งค่า timezone (optional)
RUN apk --no-cache add tzdata

# กำหนด working directory
WORKDIR /root/

# คัดลอก binary จาก builder stage
COPY --from=builder /app/app .

# Expose port ที่ใช้ (ถ้ามี)
EXPOSE 8080

# รันแอปพลิเคชัน
CMD ["./app"]
