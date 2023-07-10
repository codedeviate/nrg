// INFO: Performing a git status on all projects that is tagged at git
var oldPath = cwd();

for (const projectsKey in projects) {
    if (projects.hasOwnProperty(projectsKey)) {
        const project = projects[projectsKey];
        if(project.IsGit) {
            setbold();
            setyellow()
            println("-----------------------------------------------------------------");
            println("Project: " + project.Name);
            println("-----------------------------------------------------------------");
            setnormal();
            use(projectsKey);
            call("git", "status");
        } else {
            setbold();
            setyellow();
            println("-----------------------------------------------------------------");
            println("Project: " + project.Name + " is not a git project");
            println("-----------------------------------------------------------------");
            setnormal();
        }
    }
}
cd(oldPath);
