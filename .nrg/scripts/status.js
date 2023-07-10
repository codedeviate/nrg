// INFO: Performing a git status on the active project
var oldPath = "";
if(arguments.length > 0) {
    if(arguments[0].substring(0, 1) == "@") {
        oldPath = cwd();
        println("Old path: " + oldPath);
        const project = arguments[0];
        println("Changing to project: " + project);
        use(project);
        arguments.shift();
    }
}
call("git", "status");
if(oldPath != "") {
    cd(oldPath);
}
