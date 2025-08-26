package hover

import (
	"strings"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

func handleValueVolume(info HoverInformation) *protocol.Hover {
	valueOffset := protocol.UInteger(strings.Index(info.Line, "=")) + 1

	paths := strings.Split(info.value, ":")

	// Volume=/tmp/app:/app:ro,z
	//        ^------^ there we are
	if info.CharacterPosition >= valueOffset && info.CharacterPosition <= uint32(len(paths[0]))+valueOffset {
		msg := []string{
			"If a volume source is specified, it must be a path on the host or the name of a named volume. Host paths are allowed to be absolute or relative; relative paths are resolved relative to the directory Podman is run in. If the source does not exist, Podman returns an error. Users must pre-create the source files or directories.",
			"",
			"Any source that does not begin with a `.` or `/` is treated as the name of a named volume. If a volume with that name does not exist, it is created. Volumes created with names are not anonymous, and they are not removed by the `--rm` option and the podman rm `--volumes` command.",
		}
		return &protocol.Hover{
			Contents: protocol.MarkupContent{
				Kind:  protocol.MarkupKindMarkdown,
				Value: "**Host directory or source volume**\n\n" + strings.Join(msg, "\n"),
			},
			Range: &protocol.Range{
				Start: protocol.Position{Line: info.LineNumber, Character: valueOffset},
				End:   protocol.Position{Line: info.LineNumber, Character: valueOffset + uint32(len(paths[0]))},
			},
		}
	}

	if len(paths) < 2 {
		return nil
	}

	// Volume=/tmp/app:/app:ro,z
	//                 ^--^ there we are
	valueOffset += uint32(len(paths[0])) // +1 is the ':' character
	if info.CharacterPosition > valueOffset && info.CharacterPosition <= uint32(len(paths[1])+1)+valueOffset {
		return &protocol.Hover{
			Contents: protocol.MarkupContent{
				Kind:  protocol.MarkupKindMarkdown,
				Value: "**Container directory**\n\nThe container-dir must be an absolute path such as `/src/docs`. The volume is mounted into the container at this directory.",
			},
			Range: &protocol.Range{
				Start: protocol.Position{Line: info.LineNumber, Character: valueOffset + 1},
				End:   protocol.Position{Line: info.LineNumber, Character: valueOffset + uint32(len(paths[1])+1)},
			},
		}
	}

	if len(paths) < 3 {
		return nil
	}

	// Volume=/tmp/app:/app:ro,z
	//                      ^--^ there we are
	valueOffset += uint32(1 + len(paths[1]))
	flags := strings.SplitSeq(paths[2], ",")
	for flag := range flags {
		if info.CharacterPosition > valueOffset && info.CharacterPosition < uint32(len(flag)+1)+valueOffset {
			msg := getVolumeFlagDescription(flag)
			return &protocol.Hover{
				Contents: protocol.MarkupContent{
					Kind:  protocol.MarkupKindMarkdown,
					Value: "**Flag: " + flag + "**\n\n" + strings.Join(msg, "\n"),
				},
				Range: &protocol.Range{
					Start: protocol.Position{Line: info.LineNumber, Character: valueOffset + 1},
					End:   protocol.Position{Line: info.LineNumber, Character: valueOffset + uint32(len(flag)+1)},
				},
			}
		}
		valueOffset += uint32(len(flag) + 1)
	}

	return nil
}

func getVolumeFlagDescription(flag string) []string {
	if flag == "rw" || flag == "ro" {
		return []string{"Add `:ro` or `:rw` option to mount a volume in read-only or read-write mode, respectively. By default, the volumes are mounted read-write."}
	}

	if flag == "U" {
		return []string{
			"The `:U` suffix tells Podman to use the correct host UID and GID based on the UID and GID within the container, to change recursively the owner and group of the source volume. Chowning walks the file system under the volume and changes the UID/GID on each file.",
			"",
			"If the volume has thousands of inodes, this process takes a long time, delaying the start of the",
		}
	}

	if flag == "z" {
		return []string{
			"To change a label in the container context, add either of two suffixes :z or :Z to the volume mount. These suffixes tell Podman to relabel file objects on the shared volumes.",
			"",
			"The `z` option tells Podman that two or more containers share the volume content. As a result, Podman labels the content with a shared content label. Shared volume labels allow all containers to read/write content.",
			"",
			"Note: all containers within a pod share the same SELinux label. This means all containers within said pod can read/write volumes shared into the container created with the :Z on any one of the containers. Relabeling walks the file system under the volume and changes the label on each file; if the volume has thousands of inodes, this process takes a long time, delaying the start of the container. If the volume was previously relabeled with the z option, Podman is optimized to not relabel a second time. If files are moved into the volume, then the labels can be manually changed with the chcon -Rt container_file_t PATH command.",
			"",
			"Note: Do not relabel system files and directories. Relabeling system content might cause other confined services on the machine to fail. For these types of containers we recommend disabling SELinux separation. The option --security-opt label=disable disables SELinux separation for the container. For example if a user wanted to volume mount their entire home directory into a container, they need to disable SELinux separation.",
			"",
			"`$ podman run --security-opt label=disable -v $HOME:/home/user fedora touch /home/user/file`",
		}
	}

	if flag == "Z" {
		return []string{
			"To change a label in the container context, add either of two suffixes :z or :Z to the volume mount. These suffixes tell Podman to relabel file objects on the shared volumes.",
			"",
			"The `Z` option tells Podman to label the content with a private unshared label. Only the current container can use a private volume.",
			"",
			"Note: all containers within a pod share the same SELinux label. This means all containers within said pod can read/write volumes shared into the container created with the :Z on any one of the containers. Relabeling walks the file system under the volume and changes the label on each file; if the volume has thousands of inodes, this process takes a long time, delaying the start of the container. If the volume was previously relabeled with the z option, Podman is optimized to not relabel a second time. If files are moved into the volume, then the labels can be manually changed with the chcon -Rt container_file_t PATH command.",
			"",
			"Note: Do not relabel system files and directories. Relabeling system content might cause other confined services on the machine to fail. For these types of containers we recommend disabling SELinux separation. The option --security-opt label=disable disables SELinux separation for the container. For example if a user wanted to volume mount their entire home directory into a container, they need to disable SELinux separation.",
			"",
			"`$ podman run --security-opt label=disable -v $HOME:/home/user fedora touch /home/user/file`",
		}
	}

	if flag == "O" {
		return []string{
			"The `:O` flag tells Podman to mount the directory from the host as a temporary storage using the overlay file system. The container processes can modify content within the mountpoint which is stored in the container storage in a separate directory. In overlay terms, the source directory is the lower, and the container storage directory is the upper. Modifications to the mount point are destroyed when the container finishes executing, similar to a tmpfs mount point being unmounted.",
			"",
			"Note: The `O` flag conflicts with other options listed above.",
		}
	}

	if flag == "shared" || flag == "slave" || flag == "private" {
		return []string{
			"By default, bind-mounted volumes are private. That means any mounts done inside the container are not visible on the host and vice versa. One can change this behavior by specifying a volume mount propagation property. ",
			"",
			"When a volume is shared, mounts done under that volume inside the container are visible on host and vice versa.",
			"",
			"Making a volume slave enables only one-way mount propagation: mounts done on the host under that volume are visible inside the container but not the other way around.",
		}
	}

	if flag == "rshared" || flag == "rslave" || flag == "rprivate" {
		return []string{
			"To control mount propagation property of a volume one can use the [r]shared, [r]slave, [r]private or the [r]unbindable propagation flag. Propagation property can be specified only for bind mounted volumes and not for internal volumes or named volumes. For mount propagation to work the source mount point (the mount point where source dir is mounted on) has to have the right propagation properties. For shared volumes, the source mount point has to be shared. And for slave volumes, the source mount point has to be either shared or slave.",
		}
	}

	if flag == "rbind" {
		return []string{
			"To recursively mount a volume and all of its submounts into a container, use the rbind option. By default the bind option is used, and submounts of the source directory is not mounted into the container.",
		}
	}

	if flag == "copy" || flag == "nocopy" {
		return []string{
			"Mounting the volume with a copy option tells podman to copy content from the underlying destination directory onto newly created internal volumes. The copy only happens on the initial creation of the volume. Content is not copied up when the volume is subsequently used on different containers. The copy option is ignored on bind mounts and has no effect.",
		}
	}

	if flag == "suid" || flag == "nosuid" {
		return []string{
			"Mounting volumes with the nosuid options means that SUID executables on the volume can not be used by applications to change their privilege. By default volumes are mounted with nosuid.",
			"",
			"If the host-dir is a mount point, then dev, suid, and exec options are ignored by the kernel.",
		}
	}

	if flag == "exec" || flag == "noexec" {
		return []string{
			"Mounting the volume with the noexec option means that no executables on the volume can be executed within the container.",
			"",
			"If the host-dir is a mount point, then dev, suid, and exec options are ignored by the kernel.",
		}
	}

	if flag == "dev" || flag == "nodev" {
		return []string{
			"Mounting the volume with the nodev option means that no devices on the volume can be used by processes within the container. By default volumes are mounted with nodev.",
			"",
			"If the host-dir is a mount point, then dev, suid, and exec options are ignored by the kernel.",
		}
	}

	return []string{}
}
