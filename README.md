# SafeTransfer-Backend

## Introduction

SafeTransfer is a secure file transfer application that leverages blockchain technology for authentication and the InterPlanetary File System (IPFS) for decentralized storage. It aims to provide a secure and verifiable way of transferring files, ensuring that only the intended recipient can access the transferred files. Integrating blockchain for user authentication, SafeTransfer adds an extra layer of security, making it ideal for the transfer of sensitive or confidential files.

### Features

- **Blockchain Authentication**: Utilizes the Ethereum blockchain for user authentication, ensuring a high level of security and trust.
- **Decentralized Storage**: Leverages IPFS to store files, providing a distributed and resilient storage solution.
- **End-to-End Encryption**: Encrypts files before uploading them to IPFS, ensuring data confidentiality.
- **Signature Verification**: Validates the integrity and origin of files through digital signatures, enhancing security.

### Target Audience

SafeTransfer is designed for individuals and organizations seeking a secure method to transfer files. It is especially beneficial for:
- The legal and financial sectors, which handle confidential documents.
- Researchers sharing sensitive data.
- Anyone in need of a secure file transfer solution.

## Getting Started

### Prerequisites

Before beginning, ensure you have the following installed:
- Go (version 1.15 or later) for backend development.
- Docker and Docker Compose for running IPFS locally, or the IPFS Desktop application from the official website.
- An Ethereum wallet for blockchain interactions.

### Installation Guide

For detailed setup and configuration instructions for SafeTransfer, including both the front-end and back-end components, as well as information on the Continuous Integration/Continuous Deployment (CI/CD) pipeline, please refer to our dedicated DevOps repository:

[Safe-Transfer DevOps Installation Guide](https://github.com/light3739/Safe-Transfer-DevOps)

This guide offers comprehensive steps for getting SafeTransfer operational, covering environment setup, application configuration, and deployment strategies. It aims to assist developers and system administrators in efficiently navigating the setup process.

### Documentation

In addition to the DevOps repository, the SafeTransfer-Backend repository contains a `/docs` directory, which includes:

- **Architecture Overview**
- **API Integration**

These documents are designed to provide developers with a deeper understanding of the backend's structure and functionalities, facilitating development and integration efforts.
