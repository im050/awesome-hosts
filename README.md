# Awesome Hosts

A simple GUI for manage hosts, written in Go with the [astilectron](https://github.com/asticode/go-astilectron), and [Vue](https://github.com/vuejs/vue), [ElementUI](http://element-cn.eleme.io), etc.

![avatar](https://raw.githubusercontent.com/im050/awesome-hosts/master/screenshot/awesome-hosts.png)

## Building

Assume that you had installed the Go compilers (if not, [click here](https://golang.org/doc/install) to get the Go compilers)

#### Step 1: Install awesome hosts

Run the following commands:

    $ go get -u github.com/im050/awesome-hosts

#### Step 2: Install the bundler

Run the following commands:

    go get -u github.com/asticode/go-astilectron-bundler/
    
don't forget to add `$GOPATH/bin` to your `$PATH`.

#### Step 3: bundle the app for your current environment

Run the following commands:

    $ cd $GOPATH/src/github.com/im050/awesome-hosts
    $ astilectron-bundler -v

#### Step 4: test the app

The result is in the `output/<your os>-<your arch>` folder and is waiting for you to test it!

#### Step 5: bundle the app for more environments

To bundle the app for more environments, add an `environments` key to the bundler configuration (`bundler.json`):

```json
"environments": [
  {"arch": "amd64", "os": "linux"},
  {"arch": "386", "os": "windows"}
]
```

> The installation steps is copied from [go-astilectron-demo](https://github.com/asticode/go-astilectron-demo/)
    
## Similar projects 

* [SwitchHosts!](https://github.com/oldj/SwitchHosts) an App for managing hosts file  

> Actually, I have to admit that I copy something from <SwitchHosts!>, because it is a mature product and perfect. and, the AwesomeHosts is my first project using Golang, I will improve and perfect it for learning and growing.

## Todo List
* add/edit/delete host `ok`
* add/edit/delete group `ok`
* allow that add hosts from remote file `pending`
* add google hosts from remote as one of the default group `pending`
* clear DNS cache `pending`
* add a dock menu on Mac to provide some quick operations `pending`

## License

This project is an open-source software licensed under the [MIT License](https://github.com/im050/awesome-hosts/blob/master/LICENSE).