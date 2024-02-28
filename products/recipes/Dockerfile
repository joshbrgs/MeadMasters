# Use the official Deno image
FROM hayd/alpine-deno:1.15.2

# Set the working directory
WORKDIR /app

# Copy only the necessary files
COPY deps.ts .
COPY tsconfig.json .
COPY app.ts .

# Cache the dependencies
RUN deno cache deps.ts

# Run the Deno application
CMD ["run", "--allow-net", "--allow-write", "--allow-read", "--allow-plugin", "--allow-env", "app.ts"]

