// INFO: Performing git status on all git projects found in this folder or below
include("lib/fs.js")
include("lib/strings.js")
include("lib/git.js")

print("Finding git projects...")
const paths = sortstrings(findPaths(".git"))

if(paths.length == 0) {
    println("No git projects found")
    exit(0)
}
if(paths.length > 5) {
    print("\r" + itoa(paths.length), " projects found. Are you sure you want to pull all of them? (y/n)")
    if(readyn() === false) {
        exit(0)
    }
    println("=".repeat(60))
} else {
    println("\rFound " + paths.length + " git projects   ")
    println("=".repeat(("Found " + " git projects").length + itoa(paths.length).length))
}

const basepath = pwd()
paths.forEach(path => {
    let gitpath
    if(path == ".git") {
        gitpath = pwd()
    } else {
        gitpath = path.substring(0, path.length - 5)
    }
    cd(gitpath)
    println("Pulling " + gitpath)
    response = runcmdstr("git pull");
    println(trim(response[0]), "\n")
    cd(basepath)
})

""