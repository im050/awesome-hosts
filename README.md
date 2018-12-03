# Awesome Hosts

A simple GUI for manage hosts, written in Go with the [astilectron](https://github.com/asticode/go-astilectron), and [Vue](https://github.com/vuejs/vue), [ElementUI](http://element-cn.eleme.io), etc.

![avatar](https://raw.githubusercontent.com/im050/awesome-hosts/master/screenshot/awesome-hosts.png)

## Building

Assume that you had installed the Go compilers (if not, [click here](https://golang.org/doc/install) to get the Go compilers)

#### Step 1: Install awesome hosts

    $ go get -u github.com/im050/awesome-hosts

another way

    $ cd $GOPATH/src
    $ git clone git@github.com:im050/awesome-hosts
    $ go get -u

#### Step 2: Install the bundler

    go get -u github.com/asticode/go-astilectron-bundler/
    
don't forget to add `$GOPATH/bin` to your `$PATH`.
## Todo List
* adjust host data structure `ok`
* add host `ok`
* edit host `ok`
* delete host `ok`
* enable/disable group `ok`
* add group `ok`
* edit group name `ok`
* sync unix/window hosts `working`
* add remote hosts
* google hosts into default
* clear DNS cache `pendding`
