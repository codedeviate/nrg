// INFO: Performing git status on all git projects found in this folder or below
include("lib/fs.js")
include("lib/strings.js")
include("lib/git.js")

print("Finding git projects...")
const paths = sortstrings(glob("**/.git"))

if(paths.length == 0) {
    println("No git projects found")
    exit(0)
}

println("\rFound " + paths.length + " git projects   ")
println("=".repeat(("Found " + " git projects").length + itoa(paths.length).length))

const stats = gitpaths(paths)

const labelWidth = longeststring(stats.map(element => element.name + element.branch)) + 6

// Print it all out
stats.forEach(element => {
    const s1 = element.name
    const s2 = "(" + element.branch + ") : ";
    printmidpad(s1, labelWidth, s2);
    println(element.diff)
});

""