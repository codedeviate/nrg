const mypublicip = () => {
    const response = runcmdstr("dig +short myip.opendns.com @resolver1.opendns.com");
    return trim(response[0])
}