package zkutil

import (
	"path"

	"github.com/ironzhang/golang/log"
	"github.com/samuel/go-zookeeper/zk"
)

const (
	PERM_FILE      = zk.PermAdmin | zk.PermRead | zk.PermWrite
	PERM_DIRECTORY = zk.PermAdmin | zk.PermCreate | zk.PermDelete | zk.PermRead | zk.PermWrite
)

func DefaultACLs() []zk.ACL {
	return zk.WorldACL(zk.PermAll)
}

func DefaultDirACLs() []zk.ACL {
	return zk.WorldACL(PERM_DIRECTORY)
}

func DefaultFileACLs() []zk.ACL {
	return zk.WorldACL(PERM_FILE)
}

func ErrorEqual(a, b error) bool {
	if a != nil && b != nil {
		return a.Error() == b.Error()
	}
	return a == b
}

func CreateRecursive(conn *zk.Conn, zkPath, value string, flags int, acls []zk.ACL) (createdPath string, err error) {
	createdPath, err = conn.Create(zkPath, []byte(value), int32(flags), acls)
	if ErrorEqual(err, zk.ErrNoNode) {
		dirAcls := make([]zk.ACL, len(acls))
		for i, acl := range acls {
			dirAcls[i] = acl
			dirAcls[i].Perms = PERM_DIRECTORY
		}
		_, err = CreateRecursive(conn, path.Dir(zkPath), "", flags, dirAcls)
		if err != nil && !ErrorEqual(err, zk.ErrNodeExists) {
			return "", err
		}
		createdPath, err = conn.Create(zkPath, []byte(value), int32(flags), acls)
	}
	return
}

func WatchChildren(conn *zk.Conn, zkPath string, onChange func(children []string)) {
	CreateRecursive(conn, zkPath, "", 0, DefaultDirACLs())
	for {
		children, _, ch, err := conn.ChildrenW(zkPath)
		if err != nil {
			log.Errorf("watch children path error, path:%s, err:%v", zkPath, err)
			continue
		}
		onChange(children)
		select {
		case <-ch:
		}
	}
}
