{pkgs ? import <nixpkgs> {}}:

let
	templ = builtins.getFlake("github:a-h/templ/1e176a01b3723b169a1d79ab63db39c74a1e49f8");
in
pkgs.mkShell {
	nativeBuildInputs = with pkgs; [
		go
		nodejs

		cargo-generate
		wasm-pack
		cargo
		rustup
		rustc
		nodejs

		gopls
		delve
		go-tools

		air
		tailwindcss

		templ.packages.x86_64-linux.templ
	];
}

