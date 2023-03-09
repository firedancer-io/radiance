{ pkgs ? import <nixpkgs> { }, lib ? pkgs.lib
, buildGoModule ? pkgs.buildGoModule, }:

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
