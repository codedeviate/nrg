const gitpaths = (paths) => {
    const basepath = pwd()
    const stats = []
    paths.forEach(path => {
        let gitpath = regexp_replace(path, "\.git$", "")
        if(gitpath.length > 0) {
            cd(gitpath)
        }
        response = runcmdstr("git diff --shortstat");
        diff = trimwhitespace(response[0])
        if(diff == ""){
            diff = "No change"
        }
        response = runcmdstr("git branch --show-current");
        branch = trimwhitespace(response[0])

        if(gitpath == ""){
            gitpath = "[" + GetGitRepoName() + "]"
        }

        stats.push({
            name: gitpath,
            branch: branch,
            diff: diff
        })
        cd(basepath)
    })
    return stats
}


const GetGitRepoName = () => {
    if(fileexists(".git/config")) {
        const gitconfig = parseinifile(".git/config")
        if (gitconfig && gitconfig['remote "origin"'] && gitconfig['remote "origin"']["url"]) {
            const temp = gitconfig['remote "origin"']["url"]
            if(temp.indexOf("://") >= 0) {
                const temp2 = temp.split("://")
                if (temp2.length > 1) {
                    const temp3 = temp2.pop().split("/")
                    if (temp3.length > 1) {
                        temp3.shift()
                        const temp4 = temp3.join("/")
                        return regexp_replace(temp4, ".git$", "")
                    }
                }
            } else {
                const temp2 = temp.split(":")
                if (temp2.length > 1) {
                    const temp3 = temp2.pop()
                    return regexp_replace(temp3, ".git$", "")
                }
            }
        }
    }
    return ""
}