package main

import (
	"flag"
	"github.com/fqiyou/tools/foo/system"
	"github.com/fqiyou/tools/foo/util"
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
)
var (
	// 命令行参数
	inotify_dirs  = flag.String("inotify_dirs", "/Users/chaoyang/tmp", "监听目录")
	action_shell = flag.String("action_shell", "/Users/chaoyang/tmp/a.sh", "触发执行shell脚本名称")
	path_ext = flag.String("path_ext", ".md", "监控后缀")
	logger = util.Log
)

func execScript(shell_file string) {
	status,err := system.ExecShellScript(shell_file)
	if err != nil {
		logger.Error(err)
		logger.Error(status)
	}
}

func getDirList(dirpath string) ([]string, error) {
	var dir_list []string
	dir_err := filepath.Walk(dirpath,
		func(path string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if f.IsDir() {
				dir_list = append(dir_list, path)
				return nil
			}

			return nil
		})
	return dir_list, dir_err
}


func main() {

	flag.Parse()
	logger.SetLevel(logrus.InfoLevel)

	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		logger.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {

			case event := <-watcher.Events:
				logger.Info("action:", event)
				if !(event.Op == fsnotify.Chmod) && ".md" == string(path.Ext(path.Base(event.Name))) {
					execScript(*action_shell)

					//fi, err := os.Stat(event.Name)
					//if err == nil && !fi.IsDir() {
					//	execScript(*action_shell)
					//}
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					fi, err := os.Stat(event.Name)
					if err == nil && fi.IsDir() {
						watcher.Add(event.Name);
					}
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					fi, err := os.Stat(event.Name)
					if err == nil && fi.IsDir() {
						watcher.Remove(event.Name);
					}
				}
				if event.Op&fsnotify.Rename == fsnotify.Rename {
					watcher.Remove(event.Name)
				}

			case err := <-watcher.Errors:
				logger.Error("error:", err)
			}
		}
	}()


	dirs,err := getDirList(*inotify_dirs)
	if err != nil {
		logger.Fatal(err)
	}
	for _,file_name := range dirs{
		err = watcher.Add(file_name)
		if err != nil {
			logger.Fatal(err)
		}
	}


	<-done
}
