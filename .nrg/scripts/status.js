// INFO: Performing a git status on the active project
var oldPath = "";
var wrapPath = false;

if(arguments.length > 0) {
    if(arguments.length > 1) {
        wrapPath = true;
    }
    while(arguments.length > 0) {
        const path = arguments.shift();
        oldPath = cwd();
        if(path.substring(0, 1) == "@") {
            const project = path;
            use(project);
        } else {
            cd(path);
        }
        if(wrapPath) {
            println("=".repeat(cwd().length));
            println(cwd())
            println("=".repeat(cwd().length));
        } else {
            println(cwd())
        }
        call("git", "status");
        if(oldPath != "") {
            cd(oldPath);
        }
    }
} else {
    call("git", "status");
    if(oldPath != "") {
        cd(oldPath);
    }
}
