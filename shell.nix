{ pkgs ? import <nixpkgs> {} }:
pkgs.mkShell {
  packages = with pkgs; [ go protobuf protoc-gen-go protoc-gen-go-grpc ];
}
