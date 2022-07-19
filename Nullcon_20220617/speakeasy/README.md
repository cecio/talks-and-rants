Custom `winemu.py` file to be placed in `<SPEAKEASY_DIR>/lib/python3.9/site-packages/speakeasy/windows/`. Backup the original one to restore it afterward.

Modifications are marked with `# CP FIXME`

Also, for the correct replication of the activities done in the Webinar, you should copy 4 DLLs
- advapi32.dll
- kernel32.dll
- kernelbase.dll
- ntdll.dll

in `<SPEAKEASY_DIR>/lib/python3.9/site-packages/speakeasy/winenv/decoys/x86`
