// INFO: Performing git status on all git projects found in this folder or below
include("lib/fs.js")
include("lib/strings.js")
include("lib/git.js")

print("Finding git projects...")
const paths = sortstrings(findPaths(".git"))

if(paths.length == 0) {
    println("\rNo git projects found...")
    exit(0)
}

println("\rFound " + paths.length + " git projects   ")
println("=".repeat(("Found " + " git projects").length + itoa(paths.length).length))

const basepath = pwd()
paths.forEach(path => {
    let gitpath
    if(path != ".git") {
        gitpath = path.substring(0, path.length - 5)
        cd(gitpath)
        println("Status for " + gitpath)
    } else {
        gitpath = pwd()
        println("Status for current path")
    }
    println("-".repeat(("Status for ").length + gitpath.length))
    response = runcmdstr("git status");
    println(trim(response[0]), "\n")
    cd(basepath)
})

""