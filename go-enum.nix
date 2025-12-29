{
  buildGoModule,
  fetchFromGitHub,
  lib,
}:

buildGoModule rec {
  pname = "go-enum";
  version = "0.9.2";

  src = fetchFromGitHub {
    owner = "abice";
    repo = "go-enum";
    rev = "v${version}";
    hash = "sha256-VZH7xLEDqU8N7oU3tOWVdTABEQEp2mlh1NtTg22hzco=";
  };

  vendorHash = "sha256-bqJ+KBUsJzTNqeshq3eXFImW/JYL7zmCEwcy2xQHJeE=";

  ldflags = [
    "-s"
    "-w"
    "-X main.version=v${version}"
    "-X main.builtBy=nix"
  ];

  meta = {
    description = "An enum generator for go";
    homepage = "https://github.com/abice/go-enum";
    license = lib.licenses.mit;
    mainProgram = "go-enum";
  };
}
