insert PolicyException { 
    name := <str>$0,
    counter := <int64>$1,
    last_reconciliation := <datetime>$2,
    }
unless conflict on (.name)
else (
    update PolicyException 
    set { 
        counter := .counter + 1,
        last_reconciliation := <datetime>$2,
        }
)