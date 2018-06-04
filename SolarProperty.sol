pragma solidity ^0.4.0;
contract SolarProperty {
    
    /* declaration of specialized data types */
    enum HoldingStatus {OWNED, HELD}
    enum PaymentStatus {PAID, OVERDUE}
    
    struct Holder {
        uint percentageHeld; // must be maintained that the percentageHeld for all holders sums to 100
        HoldingStatus holdingStatus;
        uint lastFullPaymentTimestamp;
        uint unpaidBalance;
    }

    struct SolarSystem {
        string name;
        uint pricePerKWH;
        Holder[] holders;
    }
    
    /* public variables */
    address approver;
    mapping(address => SolarSystem) public solarSystems;

    /* public event on the blockchain, clients notified */
    // event AddSolarSystem(string name);
    // event AddSSHolding(address holder, uint ssIndex);
    // event Payment(address payer, uint unpaidBalance); // if paid in full, unpaidBalance = 0
    // event Repo(uint ssIndex, address lateHolder); // reposessing unpaid system, hardware can listen for this

    /* runs at initialization when contract is executed */
    constructor() public {
        approver = msg.sender;
    }

    
    function addSolarSystem(string _name, uint _pricePerKWH) public {
        require(msg.sender == approver);

        Holder storage approverHolder = Holder({
            percentageHeld: 100,
            holdingStatus: HoldingStatus.HELD,
            lastFullPaymentTimestamp: now,
            unpaidBalance: 0
        });

        SolarSystem storage newSystem = SolarSystem({
            name: _name,
            pricePerKWH: _pricePerKWH
        });
        newSystem.holders[approver] = approverHolder;

        solarSystems.push(newSystem);
    }

    /* Transfer _percentTransfer perent of holding of solar system at _targetSSAddress to _to */ 
    function addSSHolding(uint _percentTransfer, uint _targetSSAddress, address _to) public {
        require((msg.sender == approver) || msg.sender == _to);

        SolarSystem storage targetSS = solarSystems[_targetSSAddress];
        Holder[] storage holders = targetSS.holders;
        require(holders[approver].percentageHeld >= _percentTransfer);

        holders[approver].percentageHeld -= _percentTransfer;
        if (holders[_to].holdingStatus == HoldingStatus.HELD) {
            holders[_to].percentageHeld += _percentTransfer;
        } else {
            holders[_to] = Holder({
                percentageHeld: _percentTransfer,
                holdingStatus: HoldingStatus.HELD,
                lastFullPaymentTimestamp: now,
                unpaidBalance: 0
            });
        }
    }

    function removeSSHolding(uint _percentTransfer, uint _targetSSAddress, address _from) public {
        require((msg.sender == approver) || (msg.sender == _from));

        SolarSystem storage targetSS = solarSystems[_targetSSAddress];
        Holder[] storage holders = targetSS.holders;
        require(holders[_from].percentageHeld >= _percentTransfer);

        holders[_from].percentageHeld -= _percentTransfer;
        holders[approver].percentageHeld += _percentTransfer;
    }


    // function energyProduced(uint _ssAddress, uint _kWhProduced) public {
    //     SolarSystem storage producingSS = solarSystems[_ssAddress];

    //     require(producingSS.currentHolder == msg.sender);

    //     producingSS.unpaidBalance += _kWhProduced*producingSS.pricePerKWH;
    //     //TODO issue Swytch token here
    // }

    // /* payment by a currnet holder for energy consumed */
    // function pay(uint _ssAddress) payable public {
    //     SolarSystem storage targetSS = solarSystems[_ssAddress];
    //     targetSS.unpaidBalance -= msg.value;
    //     emit Payment(targetSS.currentHolder, targetSS.unpaidBalance);
    // }

    // /* transfer of ownership away from currentHolder if fails to pay */
    // function repo(uint _ssAddress) public {
    //     require(msg.sender == approver);

    //     address overdueHolder = solarSystems[_ssAddress].currentHolder;
    //     emit Repo(_ssAddress, overdueHolder);
        
    //     removePanelHolder(_ssAddress);
        
    // }
}