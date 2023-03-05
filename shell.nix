{ pkgs ? import <nixpkgs> { } }:
pkgs.mkShell { packages = with pkgs; [ go_1_19 rocksdb libpcap ]; }
