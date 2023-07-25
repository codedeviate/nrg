// INFO: Performing git status on all git projects found in this folder or below
const basepath = cwd()
const stats = []
let labelWidth = 25

print("Finding git projects...")
const s = runcmdstr("find . -name '.git'")
const paths = trim(s[0]).split("\n").sort(function (a, b) {
    return a.toLowerCase().localeCompare(b.toLowerCase());
})

if(paths.length == 0) {
    println("No git projects found")
    exit(0)
}
println("\rFound " + paths.length + " git projects   ")
println("=".repeat(("Found " + " git projects").length + itoa(paths.length).length))

paths.forEach(path => {
    if(path.length < 5){
        return
    }
    const gitpath = path.substring(0, path.length - 5)
    cd(gitpath)
    response = runcmdstr("git diff --shortstat");
    diff = trimwhitespace(response[0])
    if(diff == ""){
        diff = "No change"
    }
    response = runcmdstr("git branch --show-current");
    branch = trimwhitespace(response[0])
    stats.push({
        name: gitpath,
        branch: branch,
        diff: diff
    })
    cd(basepath)
})


stats.forEach(element => {
    let width = element.name.length + element.branch.length + 6
    if(width > labelWidth){
        labelWidth = width
    }
});


// Print it all out
stats.forEach(element => {
    const s1 = element.name
    const s2 = "(" + element.branch + ") : ";
    printmidpad(s1, labelWidth, s2);
    println(element.diff)
});

""