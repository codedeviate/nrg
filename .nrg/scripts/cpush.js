// INFO: This script is a combination of adding, committing and pushing to a git repository
function lint(arguments) {
    println("Running the linter");
    return call("nrun", "lint")[1];
}

function add(arguments) {
    println("Adding all files");
    return call("git", "add", "-A")[1];
}

function commit(arguments) {
    const message = arguments.length ? arguments[0] : "Auto-commit from nrg";
    println("Committing with message: " + message);
    return call("git", "commit", "-m", message)[1];
}

function push(arguments) {
    println("Pushing to remote");
    return call("git", "push")[1];
}

function main(arguments) {
    if(arguments.length == 0) {
        println("No commit message provided, aborting");
        return 1;
    }
    if(arguments[0].substring(0, 1) == "@") {
        const project = arguments[0].substring(1, arguments[0].length);
        println("Changing to project: " + project);
        use(project);
        arguments.shift()
    }
    if(arguments.length == 0) {
        println("No commit message provided, aborting");
        return 1;
    }

    const lintResult = lint(arguments);
    if(lintResult) {
        println("Lint failed, aborting");
        return lintResult;
    }

    const addResult = add(arguments);
    if(addResult) {
        println("Add failed, aborting");
        return addResult;
    }

    const commitResult = commit(arguments);
    if(commitResult) {
        println("Commit failed, aborting");
        return commitResult;
    }

    const pushResult = push(arguments);
    if(pushResult) {
        println("Push failed, aborting");
        return pushResult;
    }

    return 0;
}

main(arguments);
