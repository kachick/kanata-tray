{
  lib,
  buildGoModule,
  makeWrapper,
  self ? { },
  ...
}:

buildGoModule (finalAttrs: {
  pname = "kanata-tray";
  version = "git";

  src = lib.cleanSource ./..;

  vendorHash = "sha256-n6x2iXzC2a79i9J4lM45/O24e2Zf9c+x7n94d7jC0k0=";

  env.CGO_ENABLED = 0;

  flags = [ "-trimpath" ];

  ldflags = [
    "-s"
    "-w"
    "-X main.buildVersion=${(finalAttrs.version)}"
    "-X main.buildHash=${finalAttrs.src.rev or self.shortRev or self.dirtyShortRev or "unknown"}"
    "-X main.buildDate=unknown"
  ];

  buildInputs = [
    makeWrapper
  ];

  postInstall = ''
    wrapProgram $out/bin/kanata-tray --set-default KANATA_TRAY_LOG_DIR /tmp --prefix PATH : $out/bin
  '';

  meta = with lib; {
    description = "Tray Icon for Kanata";
    longDescription = ''
      A simple wrapper for kanata to control it from tray icon.
      Works on Linux.
    '';
    homepage = "https://github.com/rszyma/kanata-tray";
    license = licenses.gpl3Plus;
    platforms = platforms.linux;
  };
})
