import os
import pytest

import testinfra.utils.ansible_runner

testinfra_hosts = ['woofwoof.local']


@pytest.mark.parametrize('nginx_config_name', ['woofwoof.dog', 'meowmeow.kitty', 'asmr.meowmeow.kitty'])
def test_all_nginx_configs_exist_in_available_sites(host, nginx_config_name):
    conf = host.file(f'/etc/nginx/sites-available/{nginx_config_name}')
    assert conf.exists


@pytest.mark.parametrize('nginx_config_name', ['woofwoof.dog', 'meowmeow.kitty', 'asmr.meowmeow.kitty'])
def test_all_nginx_configs_symlinked_in_enabled_sites(host, nginx_config_name):
    conf = host.file(f'/etc/nginx/sites-enabled/{nginx_config_name}')
    assert conf.is_symlink
    assert conf.linked_to == f'/etc/nginx/sites-available/{nginx_config_name}'


def test_config__woofwoof_dot_dog(host):
    domain = 'woofwoof.dog'
    conf = host.file(f'/etc/nginx/sites-available/{domain}')

    assert conf.contains('server_name  woofwoof.dog;')


def test_config__meowmeow_dot_kitty(host):
    domain = 'meowmeow.kitty'
    conf = host.file(f'/etc/nginx/sites-available/{domain}')

    assert conf.contains('server_name www.meowmeow.kitty deep.meowmeow.kitty meowmeow.kitty;')


def test_config__asmr_dot_meowmeow_dot_kitty(host):
    domain = 'asmr.meowmeow.kitty'
    conf = host.file(f'/etc/nginx/sites-available/{domain}')

    assert conf.contains('server_name  asmr.meowmeow.kitty;')
