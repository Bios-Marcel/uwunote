#define MINGW_ROOT GetEnv("MINGW_ROOT")
#define UWU_VERSION GetEnv("UWU_VERSION")

[Setup]
AppName=UwU Note
AppVersion={#UWU_VERSION}
DefaultDirName={userappdata}\uwunote
PrivilegesRequired=none

OutputBaseFilename=uwunote-setup
OutputDir=./

[Files]
; Icons
Source: "{#MINGW_ROOT}\share\icons\Adwaita\*"; DestDir: "{app}\share\icons\Adwaita"; Flags: recursesubdirs
Source: "{#MINGW_ROOT}\share\glib-2.0\schemas\*"; DestDir: "{app}\share\glib-2.0\schemas"; Flags: recursesubdirs
; Iconloaders
Source: "{#MINGW_ROOT}\lib\gdk-pixbuf-2.0\2.10.0\*"; DestDir: "{app}\lib\gdk-pixbuf-2.0\2.10.0"; Flags: recursesubdirs
Source: "{#MINGW_ROOT}\lib\gdk-pixbuf-2.0\2.10.0\*"; DestDir: "{app}\lib\gdk-pixbuf-2.0\2.10.0"; Flags: recursesubdirs
; Libraries to run the application
Source: "{#MINGW_ROOT}\bin\libatk-1.0-0.dll"; DestDir: "{app}\bin"; Flags: recursesubdirs
Source: "{#MINGW_ROOT}\bin\libbz2-1.dll"; DestDir: "{app}\bin"; Flags: recursesubdirs
Source: "{#MINGW_ROOT}\bin\libcairo-2.dll"; DestDir: "{app}\bin"; Flags: recursesubdirs
Source: "{#MINGW_ROOT}\bin\libcairo-gobject-2.dll"; DestDir: "{app}\bin"; Flags: recursesubdirs
Source: "{#MINGW_ROOT}\bin\libdatrie-1.dll"; DestDir: "{app}\bin"; Flags: recursesubdirs
Source: "{#MINGW_ROOT}\bin\libepoxy-0.dll"; DestDir: "{app}\bin"; Flags: recursesubdirs
Source: "{#MINGW_ROOT}\bin\libexpat-1.dll"; DestDir: "{app}\bin"; Flags: recursesubdirs
Source: "{#MINGW_ROOT}\bin\libffi-6.dll"; DestDir: "{app}\bin"; Flags: recursesubdirs
Source: "{#MINGW_ROOT}\bin\libfontconfig-1.dll"; DestDir: "{app}\bin"; Flags: recursesubdirs
Source: "{#MINGW_ROOT}\bin\libfreetype-6.dll"; DestDir: "{app}\bin"; Flags: recursesubdirs
Source: "{#MINGW_ROOT}\bin\libfribidi-0.dll"; DestDir: "{app}\bin"; Flags: recursesubdirs
Source: "{#MINGW_ROOT}\bin\libgcc_s_seh-1.dll"; DestDir: "{app}\bin"; Flags: recursesubdirs
Source: "{#MINGW_ROOT}\bin\libgdk-3-0.dll"; DestDir: "{app}\bin"; Flags: recursesubdirs
Source: "{#MINGW_ROOT}\bin\libgdk_pixbuf-2.0-0.dll"; DestDir: "{app}\bin"; Flags: recursesubdirs
Source: "{#MINGW_ROOT}\bin\libgio-2.0-0.dll"; DestDir: "{app}\bin"; Flags: recursesubdirs
Source: "{#MINGW_ROOT}\bin\libglib-2.0-0.dll"; DestDir: "{app}\bin"; Flags: recursesubdirs
Source: "{#MINGW_ROOT}\bin\libgmodule-2.0-0.dll"; DestDir: "{app}\bin"; Flags: recursesubdirs
Source: "{#MINGW_ROOT}\bin\libgobject-2.0-0.dll"; DestDir: "{app}\bin"; Flags: recursesubdirs
Source: "{#MINGW_ROOT}\bin\libgraphite2.dll"; DestDir: "{app}\bin"; Flags: recursesubdirs
Source: "{#MINGW_ROOT}\bin\libgtk-3-0.dll"; DestDir: "{app}\bin"; Flags: recursesubdirs
Source: "{#MINGW_ROOT}\bin\libharfbuzz-0.dll"; DestDir: "{app}\bin"; Flags: recursesubdirs
Source: "{#MINGW_ROOT}\bin\libiconv-2.dll"; DestDir: "{app}\bin"; Flags: recursesubdirs
Source: "{#MINGW_ROOT}\bin\libintl-8.dll"; DestDir: "{app}\bin"; Flags: recursesubdirs
Source: "{#MINGW_ROOT}\bin\libjpeg-8.dll"; DestDir: "{app}\bin"; Flags: recursesubdirs
Source: "{#MINGW_ROOT}\bin\libpango-1.0-0.dll"; DestDir: "{app}\bin"; Flags: recursesubdirs
Source: "{#MINGW_ROOT}\bin\libpangocairo-1.0-0.dll"; DestDir: "{app}\bin"; Flags: recursesubdirs
Source: "{#MINGW_ROOT}\bin\libpangoft2-1.0-0.dll"; DestDir: "{app}\bin"; Flags: recursesubdirs
Source: "{#MINGW_ROOT}\bin\libpangowin32-1.0-0.dll"; DestDir: "{app}\bin"; Flags: recursesubdirs
Source: "{#MINGW_ROOT}\bin\libpcre-1.dll"; DestDir: "{app}\bin"; Flags: recursesubdirs
Source: "{#MINGW_ROOT}\bin\libpixman-1-0.dll"; DestDir: "{app}\bin"; Flags: recursesubdirs
Source: "{#MINGW_ROOT}\bin\libpng16-16.dll"; DestDir: "{app}\bin"; Flags: recursesubdirs
Source: "{#MINGW_ROOT}\bin\libstdc++-6.dll"; DestDir: "{app}\bin"; Flags: recursesubdirs
Source: "{#MINGW_ROOT}\bin\libthai-0.dll"; DestDir: "{app}\bin"; Flags: recursesubdirs
Source: "{#MINGW_ROOT}\bin\libwinpthread-1.dll"; DestDir: "{app}\bin"; Flags: recursesubdirs
Source: "{#MINGW_ROOT}\bin\zlib1.dll"; DestDir: "{app}\bin"; Flags: recursesubdirs
; Application executable
Source: "uwunote.exe"; DestDir: "{app}\bin"

[Tasks]
Name: RunAfterInstallation; Description: Run application after installation
Name: RunOnStartup; Description: Run application after starting your computer

[Icons]
Name: "{userstartup}\UwU Note"; Filename: "{app}\bin\uwunote.exe"; Tasks: RunOnStartup

[Run]
Filename: {app}\bin\uwunote.exe; WorkingDir: "{app}\bin"; Flags: postinstall nowait; Tasks: RunAfterInstallation