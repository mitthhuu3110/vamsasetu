# Build stage left out for simplicity; adjust as needed
FROM eclipse-temurin:17-jdk-alpine
WORKDIR /app
COPY target/*.jar backend.jar
EXPOSE 8080
ENTRYPOINT ["java","-jar","backend.jar"]
