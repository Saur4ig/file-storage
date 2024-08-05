# Folder Size Calculation Strategy

## Storing vs. Calculating Folder sizes

### Storing!!

I've chosen to **store folder sizes** for each folder (including all children folders) rather than calculating them every time. Here's why:

### Storing folder sizes

**Pros**:
- **Fast**: Precomputed sizes are quick to retrieve, important for busy apps.
- **Scalable**: Less work for the database since sizes aren't calculated every time.
- **Real-Time**: Updates on every file/folder change keep sizes accurate without recalculating.
- **Simple reads**: Easy to get the size without complex calculations.

**Cons**:
- **Complex writes**: Need to update sizes on every file/folder change, which is complicated.
- **Concurrency issues**: Handling multiple updates at once can be tricky.
- **Potential errors**: Bugs can lead to wrong size data. It's hard to keep data consistent if errors happen.

### Calculating folder sizes on each request

**Pros**:
- **Simple writes**: No need for complex update logic during file/folder changes.
- **Accurate**: Always gives the current size, calculated in real-time.
- **No stale data**: No risk of wrong size data due to missed updates or bugs.

**Cons**:
- **Slow**: Calculating size can be slow for large folders.
- **Scalability issues**: Adds extra work for the database, which might not handle high traffic well.
- **Complex reads**: Needs complicated queries to get size data dynamically, requiring careful indexing and optimization.

## Size calculation operations

### 1. Creating a file
- Update the folder size and all parent folder sizes in both Redis and PostgreSQL.

### 2. Creating multiple files
- Batch the size updates using a transaction mechanism.
- Update folder sizes in Redis during the upload process.
- Update folder sizes in PostgreSQL after upload is complete.

### 3. Creating a folder
- No size calculation needed for an empty folder.

### 4. Removing a single file
- Update the folder size and all parent folder sizes in both Redis and PostgreSQL.

### 5. Removing a single folder
- Recursively delete all files in the folder and update sizes in Redis and PostgreSQL.

### 6. Removing multiple folders
- Recursively delete all files and subfolders, updating sizes accordingly.
- Batch updates in Redis and PostgreSQL using a transaction mechanism.

### 7. Moving a single file
- Update sizes of the source and destination folders in both Redis and PostgreSQL.

### 8. Moving a single folder
- Recursively move all files and subfolders, updating sizes accordingly.

## Alternative ideas

### Periodic batch processing

**Description**:
- Run a scheduled job to update folder sizes regularly.

**Pros**:
- Simplifies write operations.
- Efficient for large datasets.

**Cons**:
- Folder sizes aren't updated in real-time, which might be a problem for apps needing up-to-date information.
- Needs a mechanism to ensure data consistency between updates.
