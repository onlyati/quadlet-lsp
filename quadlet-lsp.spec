Name:           quadlet-lsp
Version:        0.4.0
Release:        1%{?dist}
Summary:        Podman Quadlet Language Server

License:        GPLv3
URL:            https://github.com/onlyati/quadlet-lsp
Source0:        https://github.com/onlyati/quadlet-lsp/archive/refs/tags/v%{version}.tar.gz

BuildRequires:  golang >= 1.22
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
CGO_ENABLED=0 go build -o %{name}

%install
install -Dm0755 %{name} "%{buildroot}%{_bindir}/%{name}"

%check

%files
%license LICENSE
%doc README.md
%{_bindir}/%{name}

%changelog
* Tue Sep 16 2025 Attila Molnar <onlyati@pm.me> - 0.1.0-1
- Initial package
