# -*- coding: utf-8 -*-
from __future__ import division, absolute_import, print_function, unicode_literals
from fabric.api import *

workload = 3

env.user = 'isucon'
env.key_filename = 'id_rsa.isucon'
env.roledefs = {
    'server': ["ec2-54-64-176-75.ap-northeast-1.compute.amazonaws.com"],
}

@roles('server')
def push():
    local('GOOS=linux GOARCH=amd64 go build -o app app.go')
    local('cd prepare_script && GOOS=linux GOARCH=amd64 go build')
    sudo('supervisorctl stop isucon_go')
    put('app', 'webapp/go/app')
    put('templates', 'webapp/go')
    put('prepare_script', 'webapp/go')
    run('chmod 755 webapp/go/prepare_script/prepare_script')
    run('chmod 755 webapp/go/prepare_script/prepare.sh')
    run('chmod 755 webapp/go/app')
    sudo('supervisorctl start isucon_go')

@roles('server')
def test():
    run('grizzly start localhost:80')
    sudo('isucon3 test --workload {}'.format(workload))
    run('grizzly --executable=webapp/go/app stop localhost:80 output.prof')
    run('grizzly --format=html show output.prof > webapp/public/grizzly.html')

@roles('server')
def bench():
    run('grizzly start localhost:80')
    sudo('isucon3 benchmark --init /home/isucon/webapp/go/prepare_script/prepare.sh --workload {}'.format(workload))
    run('grizzly --executable=webapp/go/app stop localhost:80 output.prof')
    run('grizzly --format=html show output.prof > webapp/public/grizzly.html')
