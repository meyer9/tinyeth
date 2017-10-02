$(document).ready(function(){
  $('#Copy').click(function(){
    var $this = $(this);
    if (!$this.hasClass('active')) {
      $this.fadeOut('fast');
      setTimeout(function(){
        $this.addClass('active');
        $this.text('Copied');
        $this.fadeIn('fast');
      }, 150);
    }
  });

  $('#Create').click(function(e){
    if ($(this).hasClass('invalid')) {
      e.preventDefault();
    }
  })

  $("#CreateForm").keypress( function(e) {
    var chr = String.fromCharCode(e.which);
    if ("1234567890-abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ".indexOf(chr) < 0)
        return false;
  });

  $("#GenerateForm").keypress( function(e) {
    var chr = String.fromCharCode(e.which);
    if ("".indexOf(chr) < 0)
        return false;
  });

  // function addTriangleTo(target) {
  //     var dimensions = target.getClientRects()[0];
  //     var pattern = Trianglify({
  //         width: dimensions.width,
  //         height: dimensions.height
  //     });
  //     target.style['background-image'] = 'url(' + pattern.png() + ')';
  // }
  //
  // addTriangleTo(document.getElementById('Body'));
});


function copyToClipboard(element) {
  var $temp = $("<input>");
  $("body").append($temp);
  $temp.val($(element).text()).select();
  document.execCommand("copy");
  $temp.remove();
}

function validateAddress() {
  var address=document.forms["GenerateForm"]["address"].value;
  if (address==null || address=="") {
    return false;
  }
}

function validateForm() {
  validateAddress()
  var alias=document.forms["CreateForm"]["alias"].value;
  if (alias==null || alias=="") {
    $('#Create').addClass('invalid');
    $('#Create').val('Invalid');
    setTimeout(function(){
      //$('#AliasInput').val('');
      $('#Create').removeClass('invalid');
      $('#Create').val('Create');
    }, 1000)
    return false;
  } else if (alias.length < 3) {
    $('#Create').addClass('invalid');
    $('#Create').val('Too Short');
    setTimeout(function(){
      //$('#AliasInput').val('');
      $('#Create').removeClass('invalid');
      $('#Create').val('Create');
    }, 1000)
    return false;
  } else if (alias.length > 128) {
    $('#Create').addClass('invalid');
    $('#Create').val('Too Long');
    setTimeout(function(){
      //$('#AliasInput').val('');
      $('#Create').removeClass('invalid');
      $('#Create').val('Create');
    }, 1000)
    return false;
  }
}
