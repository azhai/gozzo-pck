BEGIN {
    FS = ","; OFS = "\t"
}
{
    gsub(/[ \r\n]+$/, "", $7);
    $7 = sprintf("%06d", $7)
    if (substr($6, 1, 1) != "0") {
        $6 = "0" $6
    }
    city = $3 "\t" $4 "\t" $6 "\t" $7;
    if (!(city in c)) {
        c[city] = NR;
        a[NR] = $7 "\t" city
    }
    phone = $2 "\t" $5;
    b[phone] = city
}
END {
    last = 0;
    print last,"\t\t\t\t" > "city.txt";
    n = asort(a);
    for (i = 1; i <= n; i++) {
        city = substr(a[i], 8, length(a[i])-7);
        cc[city] = i;
        print i,city >> "city.txt"
    }
    printf("") > "phone.txt"
    n = asorti(b, bkeys);
    for (i = 1; i <= n; i++) {
        phone = bkeys[i];
        city = b[phone];
        if (last != city) {
            last = city;
            print phone,cc[city] >> "phone.txt"
        }
    }
    print "2000000\t未知\t0" >> "phone.txt"
}
