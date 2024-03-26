# SafeTransfer-Backend Documentation

## Architecture Overview

The SafeTransfer backend is engineered to facilitate authentication, file management, and interactions with the IPFS network for decentralized storage. It is comprised of several pivotal components:

- **Database Management** (`/internal/db`): Handles connections and operations with the PostgreSQL database, including schema migrations for `File` and `User` models.
- **Model Definitions** (`/internal/model`): Outlines the data models for `File` and `User`, correlating them with the database structure.
- **Repository Layer** (`/internal/repository`): Provides an abstraction layer for database operations, ensuring a clear separation between the data access layer and the service logic.
- **Service Layer** (`/internal/service`): Encompasses the business logic for file uploading, downloading, user management, and authentication.
- **API Handlers** (`/internal/api`): Establishes the HTTP endpoints and request handling logic, interfacing with the service layer to process user requests.
- **Storage Integration** (`/internal/storage`): Oversees encryption, decryption, and interaction with the IPFS network for file storage and retrieval.

### Data Flow

1. **User Authentication**: Authentication is performed using Ethereum addresses and signatures, managed by the `UserService`.
2. **File Upload**: The `FileService` encrypts and uploads files to IPFS, working in conjunction with the `IPFSStorage` component.
3. **File Download**: The `DownloadService` retrieves files from IPFS, decrypts them, and serves them to the user.

## Cryptography and File Handling

The backend employs advanced cryptographic techniques for securing file transfers:

- **File Encryption and Decryption**: Utilizes AES in CTR mode for encrypting and decrypting files, ensuring data confidentiality during storage and transmission.
- **Digital Signatures**: Employs RSA for signing files, allowing for the verification of file integrity and origin.

### Encryption Process

1. **Encryption**: Files are encrypted using a generated AES key before being uploaded to IPFS.
2. **Decryption**: Upon retrieval, files are decrypted using the corresponding AES key.

### Signature Verification

- **Signing**: Files are signed using the sender's private RSA key, generating a digital signature.
- **Verification**: The digital signature is verified using the sender's public RSA key, ensuring the file's integrity and authenticity.

## API Integration

The backend provides a RESTful API for interaction with the front-end and external clients. Key endpoints include:

- **POST `/upload`**: Uploads a file, requiring authentication.
- **GET `/download/{cid}`**: Downloads a file by its CID from IPFS.
- **POST `/verifySignature`**: Verifies a user's Ethereum signature.
- **POST `/generateNonce`**: Generates a nonce for user authentication.

### Request and Response Formats

- **Upload Request**: Requires multipart form data with the file and Ethereum address.
- **Download Response**: Streams the file content with the file's SHA-256 hash included in the response headers.
- **Authentication Requests**: JSON payloads containing Ethereum addresses, nonces, and signatures.

### Error Handling

Errors are returned as JSON objects with an error message, allowing the client to handle them appropriately.
