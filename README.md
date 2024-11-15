# William_Challenge

This repo contains an ansible play to deploy a local webserver hosting a simple page "hello world." 
This page is secured so only ports 80 and 443 are accessible, and that port 80 forwards traffic to 443.

This repo also contains tests in **tests** dir to ensure the ansible play deployed ngnix properly with correct rules, and webpage is correct.

This was done on macOS with **homebrew**. the ansible play assumes the user already has brew installed on their machine.

### To install homebrew on mac, run the following in your terminal:
  ```/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"```

### The ansible play can be run with the following command:
  ```ansible-playbook -i localhost, playbook.yml --ask-become-pass```
  *note: ask-become-pass is sudo password*

### Tests in the tests directory can be executed witeh the following command:
  ```go test -v```
