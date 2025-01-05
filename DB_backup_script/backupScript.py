import os
import shutil
import time
from datetime import datetime, timedelta


DB_PATH = "route\\to\\diploma_database.db"
BACKUP_FOLDER = "route\\to\\backupFolder"
MAX_BACKUPS = 90


def create_backup():
    current_time = datetime.now().strftime("%Y-%m-%d_%Hh%Mm%Ss")
    backup_file_name = f"backup_{current_time}.db"
    backup_file_path = os.path.join(BACKUP_FOLDER, backup_file_name)
    
    try:
        shutil.copy2(DB_PATH, backup_file_path)
        print(f"Backup created: {backup_file_path}")
    except Exception as e:
        print(f"Error creating backup: {e}")


def delete_old_backups():
    files = [f for f in os.listdir(BACKUP_FOLDER) if f.startswith("backup_")]

    if len(files) > MAX_BACKUPS:
        files.sort(key=lambda f: os.path.getmtime(os.path.join(BACKUP_FOLDER, f)))
        
        oldest_backup = files[0]
        oldest_backup_path = os.path.join(BACKUP_FOLDER, oldest_backup)
        os.remove(oldest_backup_path)
        print(f"Deleted old backup: {oldest_backup}")


def main():
    create_backup()
    delete_old_backups()


if __name__ == "__main__":
    main()
