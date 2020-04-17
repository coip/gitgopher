//repo-aware iterating gopher:
//	executes tasks on a host against a freshly cloned repository.
package main

import (
	"flag"
	"os"
	"os/exec"
	"path/filepath"

	"log"
	"strings"
)

var (
	gettherepo = &exec.Cmd{
		Args: []string{"git", "clone", "https://github.com/" + ghuser + "/" + repo},
		Dir:  workdir,
	}
	tasks = []exec.Cmd{
		{
			Args: []string{"ls"},
			Dir:  repodir,
		},
		{
			Args: []string{"go", "test"},
			Dir:  repodir,
		},
		{
			Args: []string{"go", "build"},
			Dir:  repodir,
		},
	}
)

//----------------------------------------

const (
	workdir = "work"
	ghuser  = "coip"
	repo    = "moneypenny"
)

var (
	repodir = workdir + string(filepath.Separator) + repo
)

//runtime ~lightswitches
var (
	force     = flag.Bool("f", false, "will still make an attempt")
	ephemeral = flag.Bool("e", true, "cleans up workdir afterwards")
	verbose   = flag.Bool("v", false, "exposes some ~wiring for observability")
	/*
		workdir = flag.String("dir", "work", "directory to clone in. (*tasks occur in repodir, not in workdir)")
		repo    = flag.String("r", "coip/moneypenny", "github repo to run tasks against") //mitigate user and repo as seperate inputs, take one in and split it?
	*/
)

func init() {
	flag.Parse()

	//ensure workdir
	if _, err := os.Stat(workdir); os.IsNotExist(err) {
		os.Mkdir(workdir, 0700)
	}

	//remainder of init is for ensuring usrbin host dependencies are met.
	var u = make(map[string][]int)

	if *verbose {
		defer log.Printf("~$(which tasks...) => %#v", u)
	}

	for taskID, task := range tasks {
		which := task.Args[0]
		if u[which] == nil {
			//*tradeoff for current profile of:
			// [1 ~`which`/`exec.LookPath` call per unique bin]
			//is seen w/ the gophers ~deduplicative overhead here + var u.
			//added gopherload for the sake of the kernel?
			u[which] = make([]int, 0)
		}
		u[which] = append(u[which], taskID)
	}

	if srcgetter, err := exec.LookPath(gettherepo.Args[0]); err != nil {
		log.Fatal(err)
	} else {
		gettherepo.Path = srcgetter
	}

	//task{binreq,id}, assign task path to path of usrbin we find:
	for tbinreq, tid := range u {
		binpath, err := exec.LookPath(tbinreq)
		if err != nil {
			log.Fatalf("task[%d] appears infeasible: which %s?", tid, tbinreq)
		}
		//assign the path for each task that depends on the current binreq
		for _, taskID := range u[tbinreq] {
			tasks[taskID].Path = binpath
		}
	}

}

func isExistsErr(b []byte) bool {
	return strings.Contains(string(b), "already exists")
}

func main() {

	if out, err := gettherepo.CombinedOutput(); err != nil {
		if isExistsErr(out) && *ephemeral {
			if *force {
				os.RemoveAll(repodir)
				//and rerun, passing std{out,err} up to caller
				c := exec.Command(os.Args[0], os.Args[1:]...)
				c.Stdout = os.Stdout
				c.Stderr = os.Stderr
				if err := c.Run(); err != nil {
					log.Fatal(err)
				}
				os.Exit(0)
			} else {
				log.Fatal(repodir + " exists but unsure if you need those files... try -f if you insist")
			}
		}
		log.Fatal(string(out))
	}

	if *ephemeral {
		defer os.RemoveAll(repodir)
	}

	for tID, t := range tasks {
		tOut, err := t.CombinedOutput()
		if err != nil {
			log.Fatalf("%d. failed: output:[%s] and err:[%+v]", tID, tOut, err)
		} else {
			log.Printf("|| %d.) %v \n%s-----", tID, t.Args, tOut)
		}
	}

}
