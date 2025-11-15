Name:           quadlet-lsp
Version:        0.6.0
Release:        0.6.0
Summary:        Podman Quadlet Language Server

License:        GPLv3
URL:            https://github.com/onlyati/quadlet-lsp
Source0:        https://github.com/onlyati/quadlet-lsp/archive/refs/tags/v%{version}.tar.gz

BuildRequires:  golang >= 1.24.4
BuildRequires:  git

%description
This is an implementation of language server for Podman Quadlet files.

Following features are currently available:

- Code completion
  - Provide static completion based on Podman Quadlet documentation
  - Query images, volumes, networks, pods, and so on, and provide completion based on real configuration
- Hover menu
- Implemented "go definition" and "go references" functions
- Provide syntax checking
- Execute built-in commands

%prep
%autosetup -n %{name}-%{version}

%build
CGO_ENABLED=0 go build -mod=vendor -o %{name}

%install
install -Dm0755 %{name} "%{buildroot}%{_bindir}/%{name}"

%check

%files
%license LICENSE
%doc README.md
%{_bindir}/%{name}

%changelog
* Tue Nov 11 2025 Attila Molnar <onlyati@pm.me>
- Release v0.6.0 language server
* Tue Sep 25 2025 Attila Molnar <onlyati@pm.me>
- Release v0.5.0 language server
* Tue Sep 16 2025 Attila Molnar <onlyati@pm.me> 
- Initial package with v0.4.0 language server
