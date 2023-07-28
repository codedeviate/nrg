const gitpaths = (paths) => {
    const basepath = pwd()
    const stats = []
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
    return stats
}