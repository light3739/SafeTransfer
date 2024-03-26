// SPDX-License-Identifier: MIT
pragma solidity 0.8.18;

contract FileRegistry {
    // Event declaration for file registration
    event FileRegistered(address indexed owner, string cid, string fileHash);

    // Struct to hold file details
    struct FileDetails {
        address owner;
        string fileHash;
        uint256 timestamp;
    }

    // Mapping from CID to file hash
    mapping(string => string) public cidToFileHash;

    // Mapping from file hash to FileDetails
    mapping(string => FileDetails) public fileDetails;

    /**
     * @dev Registers a file using its CID and file hash.
     * Emits a FileRegistered event upon success.
     * @param cid The content identifier (CID) of the file.
     * @param fileHash The unique hash of the file.
     */
    function registerFile(string calldata cid, string calldata fileHash) external {
        require(bytes(cid).length > 0, "CID is required");
        require(bytes(fileHash).length > 0, "File hash is required");
        require(bytes(cidToFileHash[cid]).length == 0, "CID already registered");
        require(fileDetails[fileHash].timestamp == 0, "File hash already registered");

        cidToFileHash[cid] = fileHash;
        fileDetails[fileHash] = FileDetails({
            owner: msg.sender,
            fileHash: fileHash,
            timestamp: block.timestamp
        });

        emit FileRegistered(msg.sender, cid, fileHash);
    }

    /**
     * @dev Retrieves the file hash associated with a given CID.
     * @param cid The content identifier (CID) of the file.
     * @return The file hash associated with the CID.
     */
    function getFileHash(string calldata cid) external view returns (string memory) {
        require(bytes(cidToFileHash[cid]).length > 0, "CID not registered");
        return cidToFileHash[cid];
    }

    /**
     * @dev Retrieves the details of a registered file using its file hash.
     * @param fileHash The unique hash of the file.
     * @return The details of the file including owner, file hash, and timestamp.
     */
    function getFileDetails(string calldata fileHash) external view returns (FileDetails memory) {
        require(fileDetails[fileHash].timestamp != 0, "File not registered");
        return fileDetails[fileHash];
    }
}
