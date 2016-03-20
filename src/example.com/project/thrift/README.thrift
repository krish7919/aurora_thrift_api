===========================
::Install thrift compiler::
===========================
1. curl -sSL http://apache.org/dist/thrift/KEYS | gpg --import -
2. gpg --export --armor 66B778F9 | sudo apt-key add -
3. sudo bash -c "echo 'deb http://www.apache.org/dist/thrift/debian 0.9.3 main' > /etc/apt/sources.list.d/thrift.list"
4. thrift --version => Thrift version 0.9.3

===============
::Compile IDL::
===============
thrift --gen go api.thrift

