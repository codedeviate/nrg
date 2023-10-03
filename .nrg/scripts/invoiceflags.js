// INFO: Listing all Qvickly invoice flags
silent = true;

function convertIntToFlag(flag) {
    if(GetScreenWidth() >= 108) {
        var fakeBin = "";
        for (var i = flag; i < 64; i++) {
            fakeBin += "0";
        }
        fakeBin += "1";
        for (var i = 0; i < flag; i++) {
            fakeBin += "0";
        }
        return sprintf("%2d, %s", flag, fakeBin);
    }
    return sprintf("%2d", flag);
}
function checkBit(text, bit, flags) {
    if(flags & (1 << bit)) {
        return text;
    }
    return "";
}

const flagList = [
    {text: "SENT                                 ", value: 0},
    {text: "PAPER                                ", value: 1},
    {text: "MAIL2CUSTOMER                        ", value: 2},
    {text: "EBREV                                ", value: 3},
    {text: "FACTORING                            ", value: 4},
    {text: "HANDLING                             ", value: 5},
    {text: "SENTBYDISTRIBUTOR                    ", value: 6},
    {text: "CARD                                 ", value: 7},
    {text: "BANK                                 ", value: 8},
    {text: "CREDITCREATED                        ", value: 9},
    {text: "PARTPAYMENT                          ", value: 10},
    {text: "CASH                                 ", value: 11},
    {text: "CHECKOUT                             ", value: 12},
    {text: "SIMPLE                               ", value: 13},
    {text: "RECURRING_SOURCE                     ", value: 14},
    {text: "RECURRING_DEBITS                     ", value: 15},
    {text: "AUTOCREDIT                           ", value: 16},
    {text: "EXPORTED2DISTRIBUTOR                 ", value: 17},
    {text: "EFAKTURAISAVAILABLE                  ", value: 18},
    {text: "CHATUNREAD                           ", value: 19},
    {text: "QUEUED2SEND                          ", value: 20},
    {text: "SENT2INKASSO                         ", value: 21},
    {text: "FAILED2SEND                          ", value: 22},
    {text: "SENT2STRALFORS                       ", value: 23},
    {text: "SENT2TIETO                           ", value: 24},
    {text: "SWISH                                ", value: 25},
    {text: "AVRAKNAT                             ", value: 26},
    {text: "REGRESS                              ", value: 27},
    {text: "MONTHLYREPORTED                      ", value: 28},
    {text: "MULPURINFORMED                       ", value: 29},
    {text: "QUEUED2CREATEPDF                     ", value: 30},
    {text: "CREDITCHECKPENDING                   ", value: 31},
    {text: "BILLMATEAGREMENT                     ", value: 32},
    {text: "BILLMATEFEE                          ", value: 33},
    {text: "QUEUED2AUTOCANCEL                    ", value: 34},
    {text: "BILLMATEFEECREDITTOBECHECKED         ", value: 35},
    {text: "ISINDISTRIBUTORSLIST                 ", value: 36},
    {text: "QUEUEDPENDING2SEND                   ", value: 37},
    {text: "REVIEW                               ", value: 39},
    {text: "STATUS_DUE                           ", value: 40},
    {text: "STATUS_COLLECTION                    ", value: 41},
    {text: "ISBOUGHT                             ", value: 42},
    {text: "ISREGRESSABLE                        ", value: 43},
    {text: "CREDIFLOW_EFAKTURA                   ", value: 44},
    {text: "CREDIFLOW_EFAKTURA_QUEUED            ", value: 45},
    {text: "CREDIFLOW_EFAKTURA_SENT              ", value: 46},
    {text: "INTEGRATION_ISEXPORTED               ", value: 47},
    {text: "CONFIRMED_BY_ACTIVATION              ", value: 48},
    {text: "PAYMENT_EXPORTED                     ", value: 49},
    {text: "PAUSED                               ", value: 50},
    {text: "PAYMENTFLOW2                         ", value: 51},
    {text: "QUEUED2SENDPAYMENTFLOW2              ", value: 52},
    {text: "EXTENDED_DUEDATE_MERCHANT            ", value: 53},
    {text: "EXTENDED_DUEDATE_MERCHANT_NOTINVOICED", value: 54},
    {text: "EXTENDED_DUEDATE_MERCHANT_QUEUED2SEND", value: 55},
];

// println("SENT                                 ", convertIntToFlag(0));
// println("PAPER                                ", convertIntToFlag(1));
// println("MAIL2CUSTOMER                        ", convertIntToFlag(2));
// println("EBREV                                ", convertIntToFlag(3));
// println("FACTORING                            ", convertIntToFlag(4));
// println("HANDLING                             ", convertIntToFlag(5));
// println("SENTBYDISTRIBUTOR                    ", convertIntToFlag(6));
// println("CARD                                 ", convertIntToFlag(7));
// println("BANK                                 ", convertIntToFlag(8));
// println("CREDITCREATED                        ", convertIntToFlag(9));
// println("PARTPAYMENT                          ", convertIntToFlag(10));
// println("CASH                                 ", convertIntToFlag(11));
// println("CHECKOUT                             ", convertIntToFlag(12));
// println("SIMPLE                               ", convertIntToFlag(13));
// println("RECURRING_SOURCE                     ", convertIntToFlag(14));
// println("RECURRING_DEBITS                     ", convertIntToFlag(15));
// println("AUTOCREDIT                           ", convertIntToFlag(16));
// println("EXPORTED2DISTRIBUTOR                 ", convertIntToFlag(17));
// println("EFAKTURAISAVAILABLE                  ", convertIntToFlag(18));
// println("CHATUNREAD                           ", convertIntToFlag(19));
// println("QUEUED2SEND                          ", convertIntToFlag(20));
// println("SENT2INKASSO                         ", convertIntToFlag(21));
// println("FAILED2SEND                          ", convertIntToFlag(22));
// println("SENT2STRALFORS                       ", convertIntToFlag(23));
// println("SENT2TIETO                           ", convertIntToFlag(24));
// println("SWISH                                ", convertIntToFlag(25));
// println("AVRAKNAT                             ", convertIntToFlag(26));
// println("REGRESS                              ", convertIntToFlag(27));
// println("MONTHLYREPORTED                      ", convertIntToFlag(28));
// println("MULPURINFORMED                       ", convertIntToFlag(29));
// println("QUEUED2CREATEPDF                     ", convertIntToFlag(30));
// println("CREDITCHECKPENDING                   ", convertIntToFlag(31));
// println("BILLMATEAGREMENT                     ", convertIntToFlag(32));
// println("BILLMATEFEE                          ", convertIntToFlag(33));
// println("QUEUED2AUTOCANCEL                    ", convertIntToFlag(34));
// println("BILLMATEFEECREDITTOBECHECKED         ", convertIntToFlag(35));
// println("ISINDISTRIBUTORSLIST                 ", convertIntToFlag(36));
// println("QUEUEDPENDING2SEND                   ", convertIntToFlag(37));
// println("REVIEW                               ", convertIntToFlag(39));
// println("STATUS_DUE                           ", convertIntToFlag(40));
// println("STATUS_COLLECTION                    ", convertIntToFlag(41));
// println("ISBOUGHT                             ", convertIntToFlag(42));
// println("ISREGRESSABLE                        ", convertIntToFlag(43));
// println("CREDIFLOW_EFAKTURA                   ", convertIntToFlag(44));
// println("CREDIFLOW_EFAKTURA_QUEUED            ", convertIntToFlag(45));
// println("CREDIFLOW_EFAKTURA_SENT              ", convertIntToFlag(46));
// println("INTEGRATION_ISEXPORTED               ", convertIntToFlag(47));
// println("CONFIRMED_BY_ACTIVATION              ", convertIntToFlag(48));
// println("PAYMENT_EXPORTED                     ", convertIntToFlag(49));
// println("PAUSED                               ", convertIntToFlag(50));
// println("PAYMENTFLOW2                         ", convertIntToFlag(51));
// println("QUEUED2SENDPAYMENTFLOW2              ", convertIntToFlag(52));
// println("EXTENDED_DUEDATE_MERCHANT            ", convertIntToFlag(53));
// println("EXTENDED_DUEDATE_MERCHANT_NOTINVOICED", convertIntToFlag(54));
// println("EXTENDED_DUEDATE_MERCHANT_QUEUED2SEND", convertIntToFlag(55));
let argList = [];
let switchList = {
    mode: "binary"
};
for (const arg of arguments) {
    if(arg.startsWith("--")) {
        const switchName = arg.substring(2);
        if (switchName == "bin") {
            switchList.mode = "binary";
        } else if (switchName == "binary") {
            switchList.mode = "binary";
        } else if (switchName == "dec") {
            switchList.mode = "decimal";
        } else if (switchName == "decimal") {
            switchList.mode = "decimal";
        } else if (switchName == "hex") {
            switchList.mode = "hex";
        } else {
            println("Unknown switch: " + switchName);
        }
    } else if(arg.startsWith("-")) {
        const switchName = arg.substring(1);
        if(switchName == "b") {
            switchList.mode = "binary";
        } else if(switchName == "d") {
            switchList.mode = "decimal";
        } else if(switchName == "h") {
            switchList.mode = "hex";
        } else {
            println("Unknown switch: " + switchName);
        }
    } else {
        argList.push(arg);
    }
}
if(argList.length > 0) {
    for (const arg of argList) {
        let flags;
        if(switchList.mode == "binary") {
            flags = bintoints(arg);
            println("FLAGS: " + flags.join(", "));
        } else if(switchList.mode == "decimal") {
            flags = dectoints(arg);
        } else if(switchList.mode == "hex") {
            flags = hextoints(arg);
        }
        const values = [];
        for (const flagListElement of flagList) {
            let value = flagListElement.value
            let idx = 0
            while(value > 16) {
                value -= 16
                idx++
            }
            const text = checkBit(flagListElement.text, value, flags[idx])
            if(text != "") {
                values.push(trim(text));
            }
        }
        if(values.length > 0) {
            println("FLAGS: " + values.join(", "));
        } else {
            println("FLAGS: NONE");
        }
    }
} else {
    for (const flagListElement of flagList) {
        println(flagListElement.text, convertIntToFlag(flagListElement.value));
    }
}