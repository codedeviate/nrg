//INFO: Test script
println("Hello World!");
run("test2.js");
println("Arguments", arguments);
println(arguments);
println("cwd", cwd());
println("pwd", pwd());
runcmd("grep 'ahjvrobvaorv' . -R");
println("lastErrorCode", get("lastErrorCode"));
1;