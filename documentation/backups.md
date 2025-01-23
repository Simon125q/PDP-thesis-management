Backup Handling:

1. In DB_backup_script/backupScript replace the values of `DB_PATH` and `BACKUP_FOLDER` with the paths to the `diploma_database.db` file and the folder designated for database backups, respectively.
2. In DB_backup_script/backup.sh values of `PYTHON` and `SCRIPT` with path to the python interpreter and path to `backupScript.py` accordingly.
3. Make the script executable with `chmod +x backup.sh`
4. Set up a cron job to run the script regulary:
   4.1. Open the crontab editor `crontab -e`
   4.2. Create a cron job, for example to run daily at 1AM add `0 1 * * * /path/to/backup.sh` to the end of file

