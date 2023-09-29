![image](https://github.com/v3sp4n/wifi-bruteforce/assets/57196133/f8dbd271-8513-40c4-94dd-39e5c1b12a04)

SSID name Wi-Fi

password list 
<ul>
  <li>filename.txt (to log in through the passwords in this file.)</li>
  <li>range min-max (EXAMPLE"10-50",OUTPUT PASSWORD 00000010,00000011,00000012..00000050,end.(This is how you can run multiple programs to attack a single network))</li>
  <li>empty input, the program will brute force the password from 00000000 to ~130000000</li>
</ul>

When the program will make an attack - it will create a file with the name of the wi-fi with all the passwords that have been used.
If the attack is successful, a new file will be created with the password that matched the wi-fi(RESULT_TP-Link_wifi.txt).
