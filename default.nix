{ pkgs ? import (fetchTarball
  "https://github.com/NixOS/nixpkgs/archive/4e9efd3432a9a1f50e81d70e13b38a330428bcca.tar.gz")
  { }, lib ? pkgs.lib, buildGoModule ? pkgs.buildGoModule, }:

buildGoModule rec {
  pname = "radiance";
  version = "0.0.1";
  src = ./.;

  vendorHash = "sha256-qK4NJ5JCpNBtEE47JXf2fp2vLUJQLqHZNUzDC40eQMo=";

  buildInputs = with pkgs; [ rocksdb libpcap ];

  subPackages = [ "cmd/radiance" ];

  meta = with lib; {
    description = "Solana experiments, written in Go";
    homepage = "https://github.com/firedancer-io/radiance";
    license = licenses.asl20;
    platforms = platforms.linux;
  };
}
